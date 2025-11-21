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

// func (m *SerialManager) Boot(ch <-chan []byte) {
// 	errStr := "timeout"
// 	if err := m.Write([]byte(BOOT + delimeterStr)); err != nil {
// 		m.Active = false
// 		m.ResponseModelCh <- models.WSMessage{
// 			Type:  BOOT,
// 			Error: &errStr,
// 		}
// 		return
// 	}
// 	start := time.Now()
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for time.Since(start) < ackTimeout {
// 		select {
// 		case data := <-ch:
// 			dataStr := string(data)
// 			if strings.Contains(dataStr, ACK_BOOT) {
// 				m.ResponseModelCh <- models.WSMessage{
// 					Type:   BOOT,
// 					Status: &dataStr,
// 				}
// 				return
// 			}
// 		case <-ticker.C:
// 		case <-m.ctx.Done():
// 			m.ResponseModelCh <- models.WSMessage{
// 				Type:  BOOT,
// 				Error: &errStr,
// 			}
// 			return
// 		}
// 	}
// }

// func (m *SerialManager) Status(ch <-chan []byte) {
// 	if err := m.Write([]byte(STATUS + delimeterStr)); err != nil {
// 		m.Active = false
// 		return
// 	}

// 	start := time.Now()
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for time.Since(start) < ackTimeout {
// 		select {
// 		case data := <-ch:
// 			dataStr := string(data)
// 			if strings.Contains(dataStr, BUSY) ||
// 				strings.Contains(dataStr, READY) ||
// 				strings.Contains(dataStr, STAND_BY) ||
// 				strings.Contains(dataStr, ERR) {
// 				m.ResponseModelCh <- models.WSMessage{
// 					Type:   BOOT,
// 					Status: &dataStr,
// 				}
// 				return
// 			}
// 		case <-ticker.C:
// 		case <-m.ctx.Done():
// 			return
// 		}
// 	}
// }

// func (m *SerialManager) Begin(data []byte, ch <-chan []byte, msg models.WSMessage) {
// 	if err := m.Write(data); err != nil {
// 		m.Active = false
// 		errStr := DEVICE_NOT_ACTIVE
// 		m.ResponseModelCh <- models.WSMessage{
// 			Type:  BEGIN,
// 			Error: &errStr,
// 		}
// 	}
// 	start := time.Now()
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for time.Since(start) < ackTimeout {
// 		select {
// 		case data := <-ch:
// 			dataStr := string(data)
// 			if strings.Contains(dataStr, END) {
// 				// Делаем фото
// 				buf, err := m.camera.TakePhoto()
// 				if err != nil || buf == nil {
// 					ok := false
// 					m.mu.Lock()
// 					m.Control = ok
// 					m.mu.Unlock()

// 					msg.Status = &dataStr
// 					m.ResponseModelCh <- msg
// 				}

// 				// Отправляем фото на проверку
// 				_, err = m.api.RequestAI(msg.Params.Seed, *buf)
// 				if err != nil {
// 					ok := false
// 					m.mu.Lock()
// 					m.Control = ok
// 					m.mu.Unlock()
// 					msg.Status = &dataStr
// 					m.ResponseModelCh <- msg
// 				}

// 				// Если процент совпадения меньше 0.5 - считаем ее неудачной
// 				var ok bool
// 				// // if response.PercentOfMatch < 0.5 {
// 				// // 	ok = false
// 				// // 	// (models.WSMessage{
// 				// // 	// 	Type: "NEED_DECISION",
// 				// // 	// 	Payload: &models.PayloadWs{
// 				// // 	// 		Reason: fmt.Sprintf("Similarity %.2f < 0.5", response.PercentOfMatch),
// 				// // 	// 		Photo:  buf.Bytes(),
// 				// // 	// 	},
// 				// // 	// })

// 				// // 	// ЖДЁМ решения оператора
// 				// // 	select {
// 				// // 	case decision := <-m.decisionCh:
// 				// // 		if decision.OK {
// 				// // 			ok := true
// 				// // 			m.mu.Lock()
// 				// // 			m.Control = ok
// 				// // 			m.mu.Unlock()
// 				// // 			// продолжаем как будто все ОК
// 				// // 		} else {
// 				// // 			ok := false
// 				// // 			m.mu.Lock()
// 				// // 			m.Control = ok
// 				// // 			m.mu.Unlock()
// 				// // 			return models.WSMessage{
// 				// // 				Type: BEGIN,
// 				// // 			}
// 				// // 		}

// 				// 	case <-time.After(2 * time.Minute):
// 				// 		return models.WSMessage{
// 				// 			Type: BEGIN,
// 				// 		}

// 				// 	case <-m.ctx.Done():
// 				// 		return models.WSMessage{
// 				// 			Type: BEGIN,
// 				// 		}
// 				// 	}
// 				// } else {
// 				// 	ok = true
// 				// }

// 				// Если все в порядке - считаем ее удачной
// 				m.mu.Lock()
// 				m.Control = ok
// 				m.mu.Unlock()

// 				var report models.Reports
// 				now := time.Now()
// 				if ok {
// 					report = models.Reports{
// 						Shift:   int64(msg.Params.Shift),
// 						Number:  int(msg.Params.Number),
// 						Receipt: int64(msg.Params.Receipt),
// 						Turn:    int(msg.Params.Amount),
// 						Success: ok,
// 						Dt:      &now,
// 					}
// 				} else {
// 					solution := ""
// 					err := "Error"
// 					report = models.Reports{
// 						Shift:    int64(msg.Params.Shift),
// 						Number:   int(msg.Params.Number),
// 						Receipt:  int64(msg.Params.Receipt),
// 						Turn:     int(msg.Params.Amount),
// 						Success:  ok,
// 						Error:    &err,
// 						Solution: &solution,
// 						Dt:       &now,
// 					}
// 				}

// 				// добавляем отчет в базу
// 				_, err = m.repo.RepRepo.AddReports(report)
// 				if err != nil {
// 					log.Println("Error inserting report:", err)
// 				}

// 				// возвращаем каретку
// 				m.Write([]byte(RETURN + delimeterStr))
// 				start = time.Now()
// 				ticker := time.NewTicker(1 * time.Second)
// 				defer ticker.Stop()

// 				for time.Since(start) < ackTimeout {
// 					select {
// 					case data := <-ch:
// 						dataStr := string(data)
// 						if strings.Contains(dataStr, RETURN) {
// 						}
// 					case <-ticker.C:
// 					case <-m.ctx.Done():
// 						return
// 					}
// 				}
// 			}
// 		case <-ticker.C:
// 		case <-m.ctx.Done():
// 			return
// 		}
// 	}
// }

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
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for time.Since(start) < timeout {
		select {
		case data := <-ch:
			dataStr := string(data)
			if fn(dataStr) {
				return
			}
		case <-ticker.C:
		case <-m.ctx.Done():
			errStr := "timeout"
			m.ResponseModelCh <- models.WSMessage{
				Type:  msg.Type,
				Error: &errStr,
			}
			return
		}
	}
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

	checkError := func(status string) bool {
		if strings.Contains(status, ERR) {
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
			_, _ = m.repo.RepRepo.UpdateReports(report)
			return true
		}
		return false
	}
	// Wait for ACK from ESP
	fn := func(status string) bool {
		if checkError(status) {
			return true
		}
		if strings.Contains(status, "ACK") {
			m.ResponseModelCh <- models.WSMessage{
				Type:   msg.Type,
				Status: &status,
			}
			return true
		}
		return false
	}
	m.Request(code, msg, ackTimeout, ch, fn)

	// Wait for END from ESP
	fn = func(status string) bool {
		if checkError(status) {
			return true
		}
		if strings.Contains(status, END) {
			m.ResponseModelCh <- models.WSMessage{
				Type:   msg.Type,
				Status: &status,
			}
			return true
		}
		return false
	}
	m.Request(code, msg, ackTimeout, ch, fn)

	// Делаем фото
	buf, err := m.camera.GetBytesFromPhoto()
	if err != nil || buf == nil {
		ok := false
		m.mu.Lock()
		m.Control = ok
		m.mu.Unlock()

		errStr := err.Error()
		msg.Error = &errStr
		m.ResponseModelCh <- msg
	}

	// Отправляем фото на проверку
	response, err := m.api.RequestAI(msg.Params.Seed, *buf)
	if err != nil {
		ok := false
		m.mu.Lock()
		m.Control = ok
		m.mu.Unlock()

		errStr := err.Error()
		msg.Error = &errStr
		m.ResponseModelCh <- msg
	}

	bytes := buf.Bytes()

	// Если процент совпадения меньше 0.5 - считаем ее неудачной
	var ok bool
	var decisionModel models.WSMessage
	if response.PercentOfMatch < 1 || response.Seed != msg.Params.Seed {
		ok = false
		var reason string
		if response.PercentOfMatch < 1 {
			reason += fmt.Sprintf("Процент правильности посадки семян %.2f.\n", response.PercentOfMatch)
		}
		if response.Seed != msg.Params.Seed {
			reason += fmt.Sprintf("Семена %s не совпадают с выбранными семенами %s.\n", response.Seed, msg.Params.Seed)
		}
		m.ResponseModelCh <- models.WSMessage{
			Type: "NEED_DECISION",
			Payload: &models.PayloadWs{
				Reason: &reason,
				Photo:  &bytes,
			},
		}

		// ЖДЁМ решения оператора
		timeout := time.After(1 * time.Minute)
		start := time.Now()
	loop:
		for time.Since(start) < 1*time.Minute {
			select {
			case decision := <-m.DecisionCh:
				if *decision.Status == "OK" {
					ok = true
				} else {
					ok = false
				}

				decisionModel = decision
				decisionModel.Error = &reason
				break loop
			case <-timeout:
				break loop
			case <-m.ctx.Done():
			}
		}
	}

	// Если все в порядке - считаем ее удачной
	var report models.Reports
	now := time.Now()
	if ok {
		report = models.Reports{
			Id:       msg.Id,
			Shift:    int64(msg.Params.Shift),
			Number:   int(msg.Params.Number),
			Receipt:  int64(msg.Params.Receipt),
			Turn:     int(msg.Params.Turn),
			Success:  ok,
			Dt:       &now,
			Solution: decisionModel.Solution,
		}
	} else {
		report = models.Reports{
			Id:       msg.Id,
			Shift:    int64(msg.Params.Shift),
			Number:   int(msg.Params.Number),
			Receipt:  int64(msg.Params.Receipt),
			Turn:     int(msg.Params.Turn),
			Success:  ok,
			Dt:       &now,
			Error:    decisionModel.Error,
			Solution: decisionModel.Solution,
		}
	}

	updOk, err := m.repo.RepRepo.UpdateReports(report)
	if err != nil || !updOk {
		log.Println("Error updating report:", err)
	}

	// возвращаем каретку
	fn = func(status string) bool {
		if checkError(status) {
			return true
		}
		if strings.Contains(status, RETURN) {
			m.ResponseModelCh <- models.WSMessage{
				Type:   msg.Type,
				Status: &status,
			}
			return true
		}
		return false
	}
	m.Request(RETURN, msg, ackTimeout, ch, fn)

	timeout := time.After(1 * time.Minute)

	for {
		select {
		case data := <-ch:
			dataStr := string(data)
			if checkError(dataStr) {
				return
			}
			if strings.Contains(dataStr, STAND_BY) {
				msg.Status = &dataStr
				var errVal *string
				if !ok {
					errVal = decisionModel.Error
				}
				msg.Error = errVal

				if msg.Payload == nil {
					msg.Payload = &models.PayloadWs{}
				}
				msg.Payload.Control = &ok
				m.ResponseModelCh <- msg
				return
			}
		case <-timeout:
			errStr := "timeout waiting for STAND_BY"
			msg.Error = &errStr
			m.ResponseModelCh <- msg
			return
		case <-m.ctx.Done():
			return
		}
	}
}
