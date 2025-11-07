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

	ctx    context.Context
	cancel context.CancelFunc

	ResponseCh chan []byte
	ErrCh      chan error
}

func NewSerialManager(port string, baud int) *SerialManager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &SerialManager{
		portName:   port,
		baud:       baud,
		ctx:        ctx,
		cancel:     cancel,
		ResponseCh: make(chan []byte, 100),
		ErrCh:      make(chan error, 10),
	}

	go m.autoReconnect()
	return m
}

func (m *SerialManager) autoReconnect() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		if m.isConnected() {
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("STM not connected. Trying to reconnect...")

		localCtx, cancel := context.WithCancel(m.ctx)
		serialInst, err := NewSerial(m.portName, m.baud, localCtx)
		if err != nil {
			cancel()
			log.Println("Failed to connect to STM:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		m.mu.Lock()
		oldSerial := m.Serial
		m.Serial = serialInst
		m.mu.Unlock()

		if oldSerial != nil {
			_ = oldSerial.Close()
		}

		go m.listenSerial(serialInst, cancel)

		log.Println("STM connected successfully!")
	}
}

func (m *SerialManager) isConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Serial != nil
}

func (m *SerialManager) Write(data []byte) {
	m.mu.RLock()
	serial := m.Serial
	m.mu.RUnlock()

	if serial == nil {
		log.Println("⚠️ STM not connected — cannot send data")
		return
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
			select {
			case m.ErrCh <- fmt.Errorf("stm disconnected: %w", err):
			default:
			}
		}
	}
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
			select {
			case m.ResponseCh <- data:
			default:
				log.Println("⚠️ ResponseCh менеджера переполнен")
			}

		case err, ok := <-s.ErrCh:
			if !ok {
				cancel()
				return
			}
			log.Println("STM COM error:", err)
			m.mu.Lock()
			if m.Serial == s {
				m.Serial = nil
			}
			m.mu.Unlock()
			select {
			case m.ErrCh <- err:
			default:
			}
			cancel()
			return

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
