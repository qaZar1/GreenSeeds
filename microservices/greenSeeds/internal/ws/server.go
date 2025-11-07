package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	Clients    map[*Client]bool
	Serial     *device.SerialManager
	WsToSerial chan []byte
	mu         sync.RWMutex
}

func NewServer(comPath string, baud int) (*Server, error) {
	serial := device.NewSerialManager(comPath, baud)

	server := &Server{
		Clients:    make(map[*Client]bool),
		WsToSerial: make(chan []byte, 10),
		Serial:     serial,
	}

	go func() {
		for {
			if serial.Serial == nil {
				continue
			}

			for msg := range serial.ResponseCh {
				log.Println("COM read full:", string(msg))
				server.broadcast(msg)
			}
		}
	}()

	// --- Отправка данных с WS в COM ---
	go func() {
		for data := range server.WsToSerial {
			server.Serial.Write(data)
		}
	}()

	return server, nil
}

func (s *Server) broadcast(data []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for client := range s.Clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

// ------------------ WS Handler ------------------
func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New client connected")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	client := NewClient(conn)
	s.mu.Lock()
	s.Clients[client] = true
	s.mu.Unlock()

	// Чтение сообщений от клиента
	go client.ListenRead(func(msg []byte) {
		msgCopy := make([]byte, len(msg))
		copy(msgCopy, msg)
		s.WsToSerial <- msgCopy
	})

	// Отправка сообщений клиенту
	go client.ListenWrite()

	// Обработка отключения
	go func() {
		<-client.done
		conn.Close()
		s.mu.Lock()
		delete(s.Clients, client)
		s.mu.Unlock()
		log.Println("Client disconnected, total clients:", len(s.Clients))
	}()
}

// ------------------ Close ------------------
func (s *Server) Close() {
	s.Serial.Close()
	log.Println("COM port closed")
	s.mu.Lock()
	for client := range s.Clients {
		close(client.Send)
		delete(s.Clients, client)
		client.Conn.Close()
	}
	s.mu.Unlock()
}
