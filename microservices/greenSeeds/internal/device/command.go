package device

import (
	"log"
	"strings"
	"time"
)

func (m *SerialManager) UserCommand(data []byte) {
	m.resetIdleTimer(90 * time.Second)
	// подписка на ответ
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	// перед выполнением задания проверяем состояние устройства
	log.Println(time.Now())
	if strings.Contains(string(data), BEGIN) {
		if err := m.Write([]byte(BOOT + delimeterStr)); err != nil {
			m.Active = false
			return
		}
		start := time.Now()
		for time.Since(start) < 5*time.Second {
			select {
			case data := <-ch:
				log.Println("ACK_BOOT", time.Now())
				if strings.Contains(string(data), ACK_BOOT) {
					continue
				}
			case <-time.After(1 * time.Second):
			case <-m.ctx.Done():
				return
			}
		}
	}

	// отправка команды и ожидание ответа
	if err := m.Write(data); err != nil {
		m.Active = false
		return
	}
	start := time.Now()
	for time.Since(start) < 20*time.Second {
		select {
		case data := <-ch:
			log.Println("ACK_BOOT", time.Now())
			if strings.Contains(string(data), ACK_BOOT) ||
				strings.Contains(string(data), WAIT) ||
				strings.Contains(string(data), RETURN) ||
				strings.Contains(string(data), BUSY) ||
				strings.Contains(string(data), READY) ||
				strings.Contains(string(data), STAND_BY) ||
				strings.Contains(string(data), ERR) {
				return
			}

			if strings.Contains(string(data), END) {
				_, err := m.camera.TakePhoto()
				_ = err
				// if err != nil {
				// 	return
				// }
				// os.WriteFile("photo.jpg", buf.Bytes(), 0644)

				// TODO: подключить ИИ
				m.Write([]byte(RETURN + delimeterStr))
				start = time.Now()
				for time.Since(start) < 5*time.Second {
					select {
					case data := <-ch:
						if strings.Contains(string(data), RETURN) {
							return
						}
					case <-time.After(1 * time.Second):
					case <-m.ctx.Done():
						return
					}
				}
				return
			}
		case <-time.After(1 * time.Second):
		case <-m.ctx.Done():
			return
		}
	}

	// если нет ответа за 20с, делаем запрос на статус
	if err := m.Write([]byte(STATUS + delimeterStr)); err != nil {
		m.Active = false
		return
	}

	start = time.Now()
	for time.Since(start) < 5*time.Second {
		select {
		case data := <-ch:
			if strings.Contains(string(data), ACK_BOOT) ||
				strings.Contains(string(data), WAIT) ||
				strings.Contains(string(data), RETURN) ||
				strings.Contains(string(data), BUSY) ||
				strings.Contains(string(data), READY) ||
				strings.Contains(string(data), STAND_BY) ||
				strings.Contains(string(data), ERR) {
				return
			}
		case <-time.After(1 * time.Second):
		case <-m.ctx.Done():
			return
		}
	}
}
