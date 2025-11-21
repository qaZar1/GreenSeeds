package device

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/api"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/rs/zerolog"
)

type SerialManager struct {
	portName string
	baud     int

	mu      sync.RWMutex
	Serial  *Serial
	camera  *camera.Camera
	repo    *repository.Repository
	api     *api.API
	Active  bool
	Control bool

	ctx    context.Context
	cancel context.CancelFunc

	ResponseCh      chan []byte
	ResponseModelCh chan models.WSMessage
	DecisionCh      chan models.WSMessage

	idleTimer *time.Timer
	idleMu    sync.Mutex

	subsMu sync.RWMutex
	subs   []chan []byte

	log zerolog.Logger
}

func NewSerialManager(
	port string,
	baud int,
	camera *camera.Camera,
	repo *repository.Repository,
	api *api.API,
	log zerolog.Logger,
) *SerialManager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &SerialManager{
		portName:        port,
		baud:            baud,
		camera:          camera,
		repo:            repo,
		api:             api,
		ctx:             ctx,
		cancel:          cancel,
		ResponseCh:      make(chan []byte, 100),
		ResponseModelCh: make(chan models.WSMessage, 100),
		DecisionCh:      make(chan models.WSMessage, 1),
		Active:          false,
		Control:         false,

		log: log,
	}

	m.idleTimer = time.NewTimer(20 * time.Second)

	go m.idleWatcher()

	return m
}

func (m *SerialManager) idleWatcher() {
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	for {
		select {
		case <-m.idleTimer.C:
			if !m.Active {
				m.resetIdleTimer(1 * time.Minute)
				m.reconnect()
				continue
			}
			boot := []byte(BOOT + delimeterStr)
			if err := m.Write(boot); err != nil {
				log.Println("ESP not connected — cannot send data")
				m.Active = false
				m.resetIdleTimer(1 * time.Second)
				continue
			}
			start := time.Now()
			for time.Since(start) < 15*time.Second {
				select {
				case data := <-ch:
					dataStr := string(data)
					if strings.Contains(dataStr, ACK_BOOT) {
						if !m.Active {
							m.Active = true
						}
						m.resetIdleTimer(1 * time.Minute)
						m.ResponseModelCh <- models.WSMessage{
							Type:   BOOT,
							Status: &dataStr,
						}
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

func (m *SerialManager) reconnect() {
	localCtx, cancel := context.WithCancel(m.ctx)
	serialInst, err := NewSerial(m.portName, m.baud, localCtx, m.log)
	if err != nil {
		cancel()
		log.Println("Failed to connect:", err)
		m.Active = false
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
	defer m.mu.RUnlock()
	if m.Serial != nil {
		m.initialHandshake()
	}
}

func (m *SerialManager) initialHandshake() bool {
	ch := m.Subscribe()
	defer m.Unsubscribe(ch)

	boot := []byte(BOOT + delimeterStr)
	if err := m.Write(boot); err != nil {
		log.Println("ESP not connected — cannot send data")
		m.Active = false
		return false
	}
	start := time.Now()
	for time.Since(start) < 20*time.Second {
		select {
		case data := <-ch:
			dataStr := string(data)
			log.Println("COM read: ", dataStr)
			if strings.Contains(dataStr, ACK_BOOT) {
				if !m.Active {
					m.Active = true
				}
				m.resetIdleTimer(1 * time.Minute)
				m.ResponseModelCh <- models.WSMessage{
					Type:   BOOT,
					Status: &dataStr,
				}
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

func (m *SerialManager) Write(data []byte) error {
	m.mu.RLock()
	serial := m.Serial
	m.mu.RUnlock()

	if serial == nil {
		log.Println("ESP not connected, cannot send data")
		return fmt.Errorf("device not connected")
	}

	err := serial.Write(data)
	if err != nil {
		log.Printf("Error sending data: %v", err)
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

			dataStr := strings.Trim(string(data), emptyLetter)

			m.subsMu.Lock()
			for _, sub := range m.subs {
				select {
				case sub <- []byte(dataStr):
				default:
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
