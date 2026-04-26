package device

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/tarm/serial"
)

type ISerial interface {
	Run() error
	Stop() error
	Write(data []byte) error
	ReadCh() <-chan []byte
	ConnCh() <-chan ConnEvent
	Close() error
}

type Serial struct {
	port         *serial.Port
	ResponseCh   chan []byte
	mu           sync.Mutex
	log          zerolog.Logger
	ConnectionCh chan ConnEvent
	config       *serial.Config
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewSerial(port string, baud int, ctx context.Context, log zerolog.Logger) ISerial {
	ctx, cancel := context.WithCancel(ctx)
	config := &serial.Config{
		Name:        port,
		Baud:        baud,
		ReadTimeout: time.Second * 2,
	}

	s := &Serial{
		ResponseCh:   make(chan []byte, 100),
		ConnectionCh: make(chan ConnEvent, 1),
		log:          log,
		config:       config,
		ctx:          ctx,
		cancel:       cancel,
	}

	return s
}

func (s *Serial) Run() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.port != nil {
		return nil
	}

	p, err := serial.OpenPort(s.config)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to open serial port")
		return err
	}

	s.port = p
	go s.listen()

	return nil
}

func (s *Serial) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.port == nil {
		return nil
	}

	err := s.port.Close()
	s.port = nil

	return err
}

func (s *Serial) Write(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.port == nil {
		return errors.New("Port is nil")
	}

	_, err := s.port.Write(data)
	return err
}

func (s *Serial) ReadCh() <-chan []byte {
	return s.ResponseCh
}

func (s *Serial) ConnCh() <-chan ConnEvent {
	return s.ConnectionCh
}

func (s *Serial) listen() {
	buf := make([]byte, 128)
	var msg []byte

	for {
		s.mu.Lock()
		port := s.port
		s.mu.Unlock()

		if port == nil {
			return
		}

		n, err := s.port.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			s.log.Error().Err(err).Msg("Failed to read from port")
			select {
			case s.ConnectionCh <- ConnEvent{
				Type: ConnEventDisconnected,
				Err:  err,
			}:
			default:
			}

			return
		}
		if n == 0 {
			continue
		}
		for _, b := range buf[:n] {
			if b == delimeter {
				data := make([]byte, len(msg))
				copy(data, msg)

				select {
				case s.ResponseCh <- data:
				default:
				}
				msg = msg[:0]
			} else {
				msg = append(msg, b)
			}
		}
	}
}

func (s *Serial) Close() error {
	s.cancel()
	return s.Stop()
}
