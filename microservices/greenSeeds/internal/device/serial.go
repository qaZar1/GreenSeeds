package device

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type Serial struct {
	port       *serial.Port
	ResponseCh chan []byte
	mu         sync.Mutex
}

func NewSerial(port string, baud int, ctx context.Context) (*Serial, error) {
	config := &serial.Config{
		Name:        port,
		Baud:        baud,
		ReadTimeout: time.Second * 2,
	}

	p, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	s := &Serial{
		port:       p,
		ResponseCh: make(chan []byte, 100),
	}

	return s, nil
}

func (s *Serial) Write(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.port == nil {
		return fmt.Errorf("port is nil")
	}
	_, err := s.port.Write(data)
	if err != nil {
		return err
	}
	log.Printf("COM write: %s", string(data))
	return nil
}

func (s *Serial) Listen(ctx context.Context) {
	buf := make([]byte, 128)
	var msg []byte

	for {
		select {
		case <-ctx.Done():
			close(s.ResponseCh)
			return
		default:
			n, err := s.port.Read(buf)
			if err != nil {
				if err == io.EOF {
					continue
				}
				return
			}
			if n == 0 {
				continue
			}
			for _, b := range buf[:n] {
				if b == delimeter {
					data := make([]byte, len(msg))
					data = append(data, msg...)
					select {
					case s.ResponseCh <- data:
					default:
						log.Println("ResponseCh overflow, message lost")
					}
					msg = msg[:0]
				} else {
					msg = append(msg, b)
				}
			}
		}
	}
}

func (s *Serial) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.port != nil {
		err := s.port.Close()
		s.port = nil
		return err
	}
	return nil
}
