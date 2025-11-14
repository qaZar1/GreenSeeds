package device

import (
	"errors"
	"fmt"
	"log"
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
		m.mu.Lock()
		m.Control = false
		m.mu.Unlock()
		sorrel := ""
		if msg.Params.ExtraMode {
			sorrel = m.buildSorrel(msg)
		}
		data := m.buildGcode(msg, sorrel)
		if err := m.Boot(ch); err != nil {
			m.Status(ch)
			return
		}
		data = append(data, []byte(delimeterStr)...)
		m.Begin(data, ch, msg)
	case STATUS:
		m.Status(ch)
	case BOOT:
		if err := m.Boot(ch); err != nil {
			m.Status(ch)
		}
	case SETSTATUS_READY:
		if err := m.SetStatusReady(ch); err != nil {
			m.Status(ch)
		}
	}
}

func (m *SerialManager) Boot(ch <-chan []byte) error {
	if err := m.Write([]byte(BOOT + delimeterStr)); err != nil {
		m.Active = false
		return err
	}
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for time.Since(start) < ackTimeout {
		select {
		case data := <-ch:
			if strings.Contains(string(data), ACK_BOOT) {
				return nil
			}
		case <-ticker.C:
		case <-m.ctx.Done():
			return errors.New("timeout")
		}
	}

	return errors.New("timeout")
}

func (m *SerialManager) Status(ch <-chan []byte) {
	if err := m.Write([]byte(STATUS + delimeterStr)); err != nil {
		m.Active = false
		return
	}

	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for time.Since(start) < ackTimeout {
		select {
		case data := <-ch:
			if strings.Contains(string(data), BUSY) ||
				strings.Contains(string(data), READY) ||
				strings.Contains(string(data), STAND_BY) ||
				strings.Contains(string(data), ERR) {
				return
			}
		case <-ticker.C:
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *SerialManager) Begin(data []byte, ch <-chan []byte, msg models.WSMessage) error {
	if err := m.Write(data); err != nil {
		m.Active = false
		return err
	}
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for time.Since(start) < ackTimeout {
		select {
		case data := <-ch:
			if strings.Contains(string(data), END) {
				// Делаем фото
				buf, err := m.camera.TakePhoto()
				if err != nil || buf == nil {
					ok := false
					m.mu.Lock()
					m.Control = ok
					m.mu.Unlock()
					return errors.New("Error taking photo")
				}

				// Отправляем фото на проверку
				response, err := m.API.RequestAI(msg.Params.Seed, *buf)
				if err != nil {
					ok := false
					m.mu.Lock()
					m.Control = ok
					m.mu.Unlock()
					return err
				}

				// Если процент совпадения меньше 0.5 - считаем ее неудачной
				if response.PercentOfMatch < 0.5 {
					ok := false
					m.mu.Lock()
					m.Control = ok
					m.mu.Unlock()
					return errors.New("Not enough similarity")
				}

				// Если все в порядке - считаем ее удачной
				ok := true
				m.mu.Lock()
				m.Control = ok
				m.mu.Unlock()

				var report models.Reports
				now := time.Now()
				if ok {
					report = models.Reports{
						Shift:   int64(msg.Params.Shift),
						Number:  int(msg.Params.Number),
						Receipt: int64(msg.Params.Receipt),
						Turn:    int(msg.Params.Amount),
						Success: ok,
						Dt:      &now,
					}
				} else {
					solution := ""
					err := "Error"
					report = models.Reports{
						Shift:    int64(msg.Params.Shift),
						Number:   int(msg.Params.Number),
						Receipt:  int64(msg.Params.Receipt),
						Turn:     int(msg.Params.Amount),
						Success:  ok,
						Error:    &err,
						Solution: &solution,
						Dt:       &now,
					}
				}

				// добавляем отчет в базу
				_, err = m.repo.RepRepo.AddReports(report)
				if err != nil {
					log.Println("Error inserting report:", err)
				}

				// возвращаем каретку
				m.Write([]byte(RETURN + delimeterStr))
				start = time.Now()
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()

				for time.Since(start) < ackTimeout {
					select {
					case data := <-ch:
						if strings.Contains(string(data), RETURN) {
							return nil
						}
					case <-ticker.C:
					case <-m.ctx.Done():
						return errors.New("timeout")
					}
				}
				return nil
			}
		case <-ticker.C:
		case <-m.ctx.Done():
			return errors.New("timeout")
		}
	}

	return errors.New("timeout")
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

func (m *SerialManager) buildGcode(msg models.WSMessage, sorrel string) []byte {
	const query = `
BEGIN %d/%d/%d
BUNKER %d
\x02%s\x03%s`

	data := fmt.Sprintf(query,
		msg.Params.Shift,
		msg.Params.Number,
		msg.Params.Amount,
		msg.Params.Bunker,
		msg.Params.Gcode,
		sorrel)
	return []byte(data)
}

func (m *SerialManager) buildSorrel(msg models.WSMessage) string {
	const query = `
\x07Sorrel\x0A%d/%d\x0A%d/%d\x0A\x0D`

	data := fmt.Sprintf(query,
		msg.Params.Shift,
		msg.Params.Number,
		msg.Params.CompletedAmount+1,
		msg.Params.RequiredAmount)
	return data
}
