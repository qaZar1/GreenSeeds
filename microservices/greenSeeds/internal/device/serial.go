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
	ErrCh      chan error
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
		ErrCh:      make(chan error, 10),
	}

	go s.Listen(ctx)

	// Проверка связи
	if err := s.bootHandshake(ctx); err != nil {
		_ = s.Close()
		return nil, err
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

func (s *Serial) Listen(ctx context.Context) {
	buf := make([]byte, 128)
	var msg []byte

	for {
		select {
		case <-ctx.Done():
			close(s.ResponseCh)
			close(s.ErrCh)
			return
		default:
			n, err := s.port.Read(buf)
			if err != nil {
				if err == io.EOF {
					continue
				}
				select {
				case s.ErrCh <- err:
				default:
				}
				return
			}
			if n == 0 {
				continue
			}
			for _, b := range buf[:n] {
				if b == 0x04 {
					data := append([]byte(nil), msg...)
					select {
					case s.ResponseCh <- data:
					default:
						log.Println("⚠️ ResponseCh переполнен, сообщение потеряно")
					}
					msg = msg[:0]
				} else {
					msg = append(msg, b)
				}
			}
		}
	}
}

func (s *Serial) bootHandshake(ctx context.Context) error {
	boot := []byte("BOOT\x04")
	start := time.Now()
	for time.Since(start) < 10*time.Second {
		if err := s.Write(boot); err != nil {
			select {
			case s.ErrCh <- err:
			default:
			}
		}
		select {
		case data := <-s.ResponseCh:
			if string(data) == "ACK BOOT" {
				log.Println("ACK BOOT получен — STM готов")
				return nil
			}
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("timeout: STM не ответил на BOOT")
}
