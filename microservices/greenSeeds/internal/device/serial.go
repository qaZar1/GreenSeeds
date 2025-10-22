package device

import "github.com/tarm/serial"

type Serial struct {
	config *serial.Config
}

func NewSerial(port string, baud int) *Serial {
	config := &serial.Config{
		Name: port,
		Baud: baud,
	}

	return &Serial{config: config}
}

func (s *Serial) Open() (*serial.Port, error) {
	p, err := serial.OpenPort(s.config)
	if err != nil {
		return nil, err
	}

	return p, nil
}
