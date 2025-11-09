package device

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type SerialManager struct {
	portName string
	baud     int

	mu     sync.RWMutex
	Serial *Serial
	Active bool

	ctx    context.Context
	cancel context.CancelFunc

	ResponseCh chan []byte

	idleTimer *time.Timer
	idleMu    sync.Mutex

	subsMu sync.RWMutex
	subs   []chan []byte
}

func NewSerialManager(port string, baud int) *SerialManager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &SerialManager{
		portName:   port,
		baud:       baud,
		ctx:        ctx,
		cancel:     cancel,
		ResponseCh: make(chan []byte, 100),
		Active:     false,
	}

	m.idleTimer = time.NewTimer(2 * time.Second)

	go m.reconnect()

	return m
}

func (m *SerialManager) reconnect() {
	for {
		if m.Active {
			time.Sleep(3 * time.Second)
			continue
		}
		localCtx, cancel := context.WithCancel(m.ctx)
		serialInst, err := NewSerial(m.portName, m.baud, localCtx)
		if err != nil {
			cancel()
			log.Println("Failed to connect to STM:", err)
			time.Sleep(3 * time.Second)
			return
		}

		m.mu.Lock()
		oldSerial := m.Serial
		m.Serial = serialInst
		m.mu.Unlock()

		if oldSerial != nil {
			_ = oldSerial.Close()
		}

		go m.Serial.Listen(localCtx)
		go m.listenSerial(serialInst, cancel)

		m.mu.RLock()
		hasSerial := m.Serial != nil
		m.mu.RUnlock()

		log.Println("Active: ", m.Active)
		if hasSerial {
			if m.initialHandshake() {
				go m.idleWatcher()
			} else {
				m.Active = false
			}
		}
	}
}

func (m *SerialManager) initialHandshake() bool {
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	boot := []byte("BOOT\x04")
	if err := m.Write(boot); err != nil {
		log.Println("⚠️ STM not connected — cannot send data")
		m.Active = false
		return false
	}
	start := time.Now()
	for time.Since(start) < 20*time.Second {
		select {
		case data := <-ch:
			if strings.Contains(string(data), "ACK BOOT") {
				if !m.Active {
					m.Active = true
				}
				m.resetIdleTimer(1 * time.Minute)
				return true
			}
		case <-time.After(1 * time.Second):
		case <-m.ctx.Done():
			return false
		}
	}

	m.Active = false
	return false
}

func (m *SerialManager) idleWatcher() {
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	for {
		select {
		case <-m.idleTimer.C:
			boot := []byte("BOOT\x04")
			if err := m.Write(boot); err != nil {
				log.Println("⚠️ STM not connected — cannot send data")
				m.Active = false
				return
			}
			start := time.Now()
			for time.Since(start) < 15*time.Second {
				select {
				case data := <-ch:
					if strings.Contains(string(data), "ACK BOOT") ||
						strings.Contains(string(data), "END") {
						if !m.Active {
							m.Active = true
						}
						m.resetIdleTimer(1 * time.Minute)
						continue
					}
				case <-time.After(1 * time.Second):
				case <-m.ctx.Done():
					return
				}
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *SerialManager) UserCommand(data []byte) {
	m.resetIdleTimer(90 * time.Second)
	// подписка на ответ
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	// перед выполнением задания проверяем состояние устройства
	if strings.Contains(string(data), "BEGIN") {
		if err := m.Write([]byte("BOOT\x04")); err != nil {
			m.Active = false
			return
		}
		start := time.Now()
		for time.Since(start) < 20*time.Second {
			select {
			case data := <-ch:
				if strings.Contains(string(data), "ACK BOOT") {
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
			if strings.Contains(string(data), "ACK") ||
				strings.Contains(string(data), "WAIT") ||
				strings.Contains(string(data), "RETURN") ||
				strings.Contains(string(data), "BUSY") ||
				strings.Contains(string(data), "READY") ||
				strings.Contains(string(data), "STAND BY") ||
				strings.Contains(string(data), "ERR") {
				return
			}
		case <-time.After(1 * time.Second):
		case <-m.ctx.Done():
			return
		}
	}

	// если нет ответа за 20с, делаем запрос на статус
	if err := m.Write([]byte("BOOT\x04")); err != nil {
		m.Active = false
		return
	}

	start = time.Now()
	for time.Since(start) < 20*time.Second {
		select {
		case data := <-ch:
			if strings.Contains(string(data), "ACK BOOT") {
				return
			}
		case <-time.After(1 * time.Second):
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *SerialManager) Write(data []byte) error {
	m.mu.RLock()
	serial := m.Serial
	m.mu.RUnlock()

	if serial == nil {
		log.Println("⚠️ STM not connected — cannot send data")
		return fmt.Errorf("device not connected")
	}

	err := serial.Write(data)
	if err != nil {
		log.Printf("⚠️ Ошибка отправки: %v", err)
		if strings.Contains(err.Error(), "device not configured") ||
			strings.Contains(err.Error(), "input/output error") {
			m.mu.Lock()
			if m.Serial == serial {
				m.Serial = nil
			}
			m.mu.Unlock()
		}

		return err
	}

	return nil
}

func (m *SerialManager) listenSerial(s *Serial, cancel context.CancelFunc) {
	for {
		select {
		case data, ok := <-s.ResponseCh:
			if !ok {
				cancel()
				m.mu.Lock()
				if m.Serial == s {
					m.Serial = nil
				}
				m.mu.Unlock()
				return
			}

			m.subsMu.Lock()
			for _, sub := range m.subs {
				select {
				case sub <- data:
				default:
					log.Println("⚠️ ResponseCh переполнен, сообщение потеряно")
				}
			}
			m.subsMu.Unlock()

		case <-m.ctx.Done():
			cancel()
			return
		}
	}
}

func (m *SerialManager) Close() {
	m.cancel()
	m.mu.Lock()
	if m.Serial != nil {
		_ = m.Serial.Close()
		m.Serial = nil
	}
	m.mu.Unlock()
}

func (m *SerialManager) resetIdleTimer(dur time.Duration) {
	m.idleMu.Lock()
	defer m.idleMu.Unlock()
	if m.idleTimer == nil {
		m.idleTimer = time.NewTimer(dur)
	}
	if !m.idleTimer.Stop() {
		select {
		case <-m.idleTimer.C:
		default:
		}
	}
	m.idleTimer.Reset(dur)
}
