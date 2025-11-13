package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	Clients     map[*Client]bool
	Serial      *device.SerialManager
	WsToSerial  chan []byte
	mu          sync.RWMutex
	isBusy      bool
	commandType string
	amount      int
}

func NewServer(comPath string, baud int, camera *camera.Camera, repo *repository.Repository) (*Server, error) {
	serial := device.NewSerialManager(comPath, baud, camera, repo)

	server := &Server{
		Clients:    make(map[*Client]bool),
		WsToSerial: make(chan []byte, 10),
		Serial:     serial,
	}

	// Чтение данных с COM
	go func() {
		ch := server.Serial.Subscribe()
		defer server.Serial.Unsubscribe(ch)

		for msg := range ch {
			log.Println("COM read:", string(msg))
			cleanedMsg := strings.Trim(string(msg), "\x00")

			var wsMsg models.WSMessage

			server.mu.Lock()
			commandType := server.commandType
			amount := server.amount
			server.mu.Unlock()

			// Берём актуальное значение control из Serial
			control := server.Serial.Control

			if cleanedMsg == "ACK BOOT" {
				wsMsg = models.WSMessage{
					Type:   "BOOT",
					Status: &cleanedMsg,
				}
			} else {
				wsMsg = models.WSMessage{
					Type:   commandType,
					Status: &cleanedMsg,
					Params: &models.Params{
						Amount: amount,
					},
					Payload: &models.PayloadWs{
						Control: control,
					},
				}
			}

			server.broadcast(wsMsg)

			server.mu.Lock()
			server.isBusy = false
			server.mu.Unlock()
		}
	}()

	// Отправка данных с WS в COM
	go func() {
		for data := range server.WsToSerial {
			// анмаршалим данные в модель
			var msg models.WSMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Println("WS to COM error:", err)
				err := "Failed to parse message"
				wsMsg := models.WSMessage{
					Type:  "ERR",
					Error: &err,
				}
				server.broadcast(wsMsg)
				continue
			}

			// устройство занято или неактивно
			if !server.Serial.Active || server.isBusy {
				err := "Device is not active"
				wsMsg := models.WSMessage{
					Type:  msg.Type,
					Error: &err,
				}
				server.broadcast(wsMsg)
				continue
			}

			// устанавливаем флаг занятости, копируем тип команды и количество
			server.mu.Lock()
			server.isBusy = true
			server.commandType = msg.Type
			if msg.Params != nil && msg.Params.Amount != 0 {
				server.amount = msg.Params.Amount
			}
			server.mu.Unlock()

			server.Serial.UserCommand(msg)
		}
	}()

	return server, nil
}

func (s *Server) broadcast(msg models.WSMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for client := range s.Clients {
		select {
		case client.Send <- msg:
		default:
		}
	}
}

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
