package device

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (m *SerialManager) UserCommand(msg models.WSMessage) {
	m.resetIdleTimer(90 * time.Second)

	// подписка на ответ
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	switch msg.Type {
	case BEGIN:
		m.Begin(msg, ch)
	case STATUS:
		fn := func(status string) bool {
			if strings.Contains(status, BUSY) ||
				strings.Contains(status, READY) ||
				strings.Contains(status, STAND_BY) ||
				strings.Contains(status, WAIT) ||
				strings.Contains(status, ERR) {
				m.ResponseModelCh <- models.WSMessage{
					Type:   STATUS,
					Status: &status,
				}
				return true
			}
			return false
		}

		m.Request(STATUS, msg, ackTimeout, ch, fn)
	case BOOT:
		fn := func(status string) bool {
			if strings.Contains(status, ACK_BOOT) {
				m.ResponseModelCh <- models.WSMessage{
					Type:   BOOT,
					Status: &status,
				}
				return true
			}
			return false
		}

		m.Request(BOOT, msg, ackTimeout, ch, fn)
	case SETSTATUS_READY:
		fn := func(status string) bool {
			if strings.Contains(status, ACK_SETSTATUS_READY) {
				m.ResponseModelCh <- models.WSMessage{
					Type:   SETSTATUS_READY,
					Status: &status,
				}
				return true
			}
			return false
		}

		m.Request(SETSTATUS_READY, msg, ackTimeout, ch, fn)
	}
}

func (m *SerialManager) SetStatusReady(ch <-chan []byte) error {
	if err := m.Write([]byte(SETSTATUS_READY + delimeterStr)); err != nil {
		m.Active = false
		return err
	}
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for time.Since(start) < ackTimeout {
		select {
		case data := <-ch:
			if strings.Contains(string(data), READY) {
				return nil
			}
		case <-ticker.C:
		case <-m.ctx.Done():
			return errors.New("timeout")
		}
	}

	return errors.New("timeout")
}

func (m *SerialManager) Request(
	write string,
	msg models.WSMessage,
	timeout time.Duration,
	ch <-chan []byte,
	fn func(string) bool,
) {
	if err := m.Write([]byte(write + delimeterStr)); err != nil {
		m.Active = false
		errStr := err.Error()
		m.ResponseModelCh <- models.WSMessage{
			Type:  msg.Type,
			Error: &errStr,
		}
		return
	}

	err := m.wait(ch, timeout, fn)
	if err != nil {
		errStr := "timeout"
		m.ResponseModelCh <- models.WSMessage{
			Type:  msg.Type,
			Error: &errStr,
		}
	}
}

func (m *SerialManager) wait(
	ch <-chan []byte,
	timeout time.Duration,
	onData func(string) bool,
) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	for {
		select {
		case data := <-ch:
			if onData(string(data)) {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *SerialManager) handleDeviceError(msg models.WSMessage, status string) bool {
	if !strings.Contains(status, ERR) {
		return false
	}

	errStr := "Device returned ERR: " + status
	msg.Error = &errStr
	msg.Status = &status
	m.ResponseModelCh <- msg

	now := time.Now()
	report := models.Reports{
		Id:      msg.Id,
		Shift:   int64(msg.Params.Shift),
		Number:  int(msg.Params.Number),
		Receipt: int64(msg.Params.Receipt),
		Turn:    int(msg.Params.Turn),
		Success: false,
		Error:   &errStr,
		Dt:      &now,
	}

	m.repo.RepRepo.UpdateReports(report)
	return true
}

func (m *SerialManager) waitDecision(
	reason string,
	photo []byte,
) (bool, models.WSMessage, error) {
	m.ResponseModelCh <- models.WSMessage{
		Type: "NEED_DECISION",
		Payload: &models.PayloadWs{
			Reason: &reason,
			Photo:  &photo,
		},
	}

	select {
	case decision := <-m.DecisionCh:
		ok := decision.Status != nil && *decision.Status == "OK"
		decision.Error = &reason
		return ok, decision, nil

	case <-time.After(time.Minute):
		return false, models.WSMessage{}, errors.New("decision timeout")

	case <-m.ctx.Done():
		return false, models.WSMessage{}, m.ctx.Err()
	}
}

func (m *SerialManager) buildReport(
	msg models.WSMessage,
	ok bool,
	decision models.WSMessage,
) models.Reports {
	now := time.Now()

	report := models.Reports{
		Id:      msg.Id,
		Shift:   int64(msg.Params.Shift),
		Number:  int(msg.Params.Number),
		Receipt: int64(msg.Params.Receipt),
		Turn:    int(msg.Params.Turn),
		Success: ok,
		Dt:      &now,
	}

	if !ok {
		report.Error = decision.Error
	}

	report.Solution = decision.Solution
	return report
}

func (m *SerialManager) Begin(msg models.WSMessage, ch <-chan []byte) {
	m.mu.Lock()
	m.Control = false
	m.mu.Unlock()
	var sorrel string
	if msg.Params.ExtraMode {
		sorrel = m.buildSorrel(msg)
	}
	code := m.buildGcode(msg, sorrel)

	check := func(status string) bool {
		return m.handleDeviceError(msg, status)
	}

	// ACK
	m.Request(code, msg, ackTimeout, ch, func(s string) bool {
		if check(s) {
			return true
		}
		if strings.Contains(s, "ACK") {
			m.ResponseModelCh <- models.WSMessage{Type: msg.Type, Status: &s}
			return true
		}
		return false
	})

	// END
	m.Request(code, msg, ackTimeout, ch, func(s string) bool {
		if check(s) {
			return true
		}
		if strings.Contains(s, END) {
			m.ResponseModelCh <- models.WSMessage{Type: msg.Type, Status: &s}
			return true
		}
		return false
	})

	isAffected, err := m.repo.PlcRepo.DecrementSeed(msg.Params.Bunker)
	if err != nil || !isAffected {
		text := "Нет семян в бункере"
		m.ResponseModelCh <- models.WSMessage{
			Type:   "ERR",
			Error:  &text,
			Params: msg.Params,
		}
		return
	}

	bunkers, err := m.repo.SeedRepo.GetSeedsWithBunkers(msg.Params.Seed)
	if err != nil {
		text := "Нет семян в бункерах"
		m.ResponseModelCh <- models.WSMessage{
			Type:   "ERR",
			Error:  &text,
			Params: msg.Params,
		}
		return
	}

	m.ResponseModelCh <- models.WSMessage{
		Type:    "BUNKERS_UPDATE",
		Bunkers: &bunkers,
	}

	buf, err := m.camera.GetBytesFromPhoto("./../../test1.jpg")
	if err != nil || buf == nil {
		errStr := err.Error()
		msg.Error = &errStr
		m.ResponseModelCh <- msg
		return
	}

	response, err := m.api.RequestAI(msg.Params.Seed, *buf)
	if err != nil {
		errStr := err.Error()
		msg.Error = &errStr
		m.ResponseModelCh <- msg
		return
	}

	ok := true
	var decision models.WSMessage
	bytes := buf.Bytes()

	if response.PercentOfMatch < 1 || response.Seed != msg.Params.Seed {
		ok = false
		var reason string

		if response.PercentOfMatch < 1 {
			reason += fmt.Sprintf("Процент посадки %.2f\n", response.PercentOfMatch)
		}
		if response.Seed != msg.Params.Seed {
			reason += fmt.Sprintf("Семена %s ≠ %s\n", response.Seed, msg.Params.Seed)
		}

		ok, decision, err = m.waitDecision(reason, bytes)
		if err != nil {
			return
		}
	}

	report := m.buildReport(msg, ok, decision)
	m.repo.RepRepo.UpdateReports(report)

	// RETURN
	m.Request(RETURN, msg, ackTimeout, ch, func(s string) bool {
		if check(s) {
			return true
		}
		return strings.Contains(s, RETURN)
	})

	// STAND_BY
	m.wait(ch, time.Minute, func(s string) bool {
		if check(s) {
			return true
		}
		if strings.Contains(s, STAND_BY) {
			msg.Status = &s
			msg.Payload = &models.PayloadWs{
				Control: &ok,
			}
			if !ok {
				msg.Error = decision.Error
			}
			m.ResponseModelCh <- msg
			return true
		}
		return false
	})
}
