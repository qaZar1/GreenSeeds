package device

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Manager struct {
	mu         sync.RWMutex
	serial     ISerial
	dispatcher *Dispatcher

	cmdMu sync.Mutex

	ctx  context.Context
	port string
	baud int
	log  zerolog.Logger

	reconnecting   bool
	status         ManagerState
	StatusChangeCh chan ManagerState

	state DeviceState

	activeSessionId string
	sessionMu       sync.RWMutex
}

const InternalSessionID = "system"

func NewManager(ctx context.Context, port string, baud int, log zerolog.Logger) *Manager {
	m := Manager{
		ctx:  ctx,
		port: port,
		baud: baud,
		log:  log,

		reconnecting: false,

		StatusChangeCh: make(chan ManagerState, 10),
	}

	serial := NewSerial(m.port, m.baud, m.ctx, m.log)
	dispatcher := NewDispatcher()
	m.serial = serial
	m.dispatcher = dispatcher
	go m.handleSerial(serial)

	if err := serial.Run(); err != nil {
		m.status = ManagerStateConnecting

		m.startReconnect()
	} else {
		m.status = ManagerStateConnected
	}

	return &m
}

func (m *Manager) startReconnect() {
	m.mu.Lock()
	if m.reconnecting {
		m.mu.Unlock()
		return
	}
	m.status = ManagerStateConnecting
	m.reconnecting = true
	m.mu.Unlock()

	go m.reconnect()
}

func (m *Manager) reconnect() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
		case <-ticker.C:
			if err := m.serial.Run(); err != nil {
				continue
			}

			m.mu.Lock()
			m.reconnecting = false
			m.status = ManagerStateConnected
			m.mu.Unlock()
			m.WriteStatusCh(ManagerStateConnected)

			return
		}
	}
}

func (m *Manager) Write(data []byte) error {
	m.mu.RLock()
	serial := m.serial
	m.mu.RUnlock()

	if serial == nil {
		return errors.New("Serial is nil")
	}

	if err := serial.Write(data); err != nil {
		m.mu.Lock()
		m.status = ManagerStateDisconnected
		m.mu.Unlock()

		serial.Stop()

		m.startReconnect()
		return err
	}

	return nil
}

func (m *Manager) handleSerial(serial ISerial) {
	readCh := serial.ReadCh()
	connCh := serial.ConnCh()

	for {
		select {
		case data, ok := <-readCh:
			if !ok {
				m.serial.Stop()
				m.helperHandle(serial)
				continue
			}

			m.dispatcher.Handle(data)
		case event := <-connCh:
			if event.Type == ConnEventDisconnected {
				m.serial.Stop()
				m.helperHandle(serial)
				continue
			}
		}
	}
}

func (m *Manager) helperHandle(serial ISerial) {
	m.mu.Lock()
	m.status = ManagerStateDisconnected
	m.mu.Unlock()

	m.WriteStatusCh(ManagerStateDisconnected)

	// сброс сессии при отключении устройства
	m.sessionMu.Lock()
	if m.activeSessionId != "" {
		m.log.Warn().Msg("session terminated due to disconnect")
		m.activeSessionId = ""
	}
	m.sessionMu.Unlock()

	m.dispatcher.FailAll(errors.New("disconnected"))

	m.startReconnect()
}

func (m *Manager) GetStatus() ManagerState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

func (m *Manager) Do(
	ctx context.Context,
	cmd []byte,
	match func([]byte) MatchResult,
	buf int,
	callerSessionId string,
) (<-chan []byte, error) {
	m.sessionMu.RLock()
	if m.activeSessionId != "" && m.activeSessionId != callerSessionId {
		m.sessionMu.RUnlock()
		return nil, errors.New("session already acquired")
	}
	m.sessionMu.RUnlock()

	m.cmdMu.Lock()
	req := &Request{
		Ch:    make(chan []byte, buf),
		Done:  make(chan error, 1),
		Match: match,
	}

	out := make(chan []byte, buf)

	go func() {
		defer close(out)
		defer m.cmdMu.Unlock()

		for {
			select {
			case <-req.Done:
				m.dispatcher.remove(req)
				return

			case <-ctx.Done():
				m.dispatcher.remove(req)
				return

			case msg, ok := <-req.Ch:
				if !ok {
					return
				}

				select {
				case out <- msg:
				case <-req.Done:
					return
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	m.dispatcher.Add(req)

	if err := m.Write(cmd); err != nil {
		m.dispatcher.remove(req)
		return nil, err
	}

	return out, nil
}

func (m *Manager) GetState() DeviceState {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.state
}

func (m *Manager) SetState(newState DeviceState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state = newState
}

func (m *Manager) WriteStatusCh(status ManagerState) {
	select {
	case m.StatusChangeCh <- status:
	default:
	}
}

func (m *Manager) parseStatus(msg string) {
	switch {
	case strings.HasPrefix(msg, MANUAL_MODE):
		m.SetState(StateManual)
	case strings.HasPrefix(msg, STAND_BY):
		m.SetState(StateStandby)
	case strings.HasPrefix(msg, READY):
		m.SetState(StateReady)
	case strings.HasPrefix(msg, BUSY):
		m.SetState(StateBusy)
	case strings.HasPrefix(msg, WAIT):
		m.SetState(StateWait)
	case strings.HasPrefix(msg, RETURN):
		m.SetState(StateReturn)
	case strings.HasPrefix(msg, ERR):
		m.SetState(StateError)
	}
}

// func (m *Manager) updateState(data []byte) {
// 	newState := parseStatus(data)

// 	if newState == StateUnknown {
// 		return
// 	}

// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	if m.state != newState {
// 		m.state = newState
// 		m.log.Info().
// 			Str("new_state", newState.String()).
// 			Msg("device state updated")
// 	}
// }
