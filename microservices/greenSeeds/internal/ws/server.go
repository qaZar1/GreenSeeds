package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/api"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	Clients    map[*Client]bool
	Serial     *device.SerialManager
	WsToSerial chan models.WSMessage
	SerialToWs chan models.WSMessage
	mu         sync.RWMutex
	isBusy     bool
	log        zerolog.Logger
}

func NewServer(
	comPath string,
	baud int,
	camera *camera.Camera,
	repo *repository.Repository,
	url string,
	log zerolog.Logger,
) (*Server, error) {
	api := api.NewAPI(url)
	serial := device.NewSerialManager(comPath, baud, camera, repo, api, log)

	server := &Server{
		Clients:    make(map[*Client]bool),
		WsToSerial: make(chan models.WSMessage, 10),
		Serial:     serial,
		log:        log,
	}

	// Чтение данных с COM
	ch := serial.Subscribe()
	defer serial.Unsubscribe(ch)

	go func() {
		for msg := range server.Serial.ResponseModelCh {
			server.broadcast(msg)
		}
	}()

	// Отправка данных с WS в COM
	go func() {
		for msg := range server.WsToSerial {
			if msg.Type == "DECISION" {
				server.Serial.DecisionCh <- msg
			}
			// устройство занято или неактивно
			if !server.Serial.Active {
				err := "Device is not active"
				wsMsg := models.WSMessage{
					Type:  msg.Type,
					Error: &err,
				}
				server.Serial.ResponseModelCh <- wsMsg
				continue
			}

			// устанавливаем флаг занятости
			server.mu.Lock()
			server.isBusy = true
			server.mu.Unlock()

			go server.Serial.UserCommand(msg)
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
		var wsMsg models.WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			log.Println("WS to COM error:", err)
			return
		}

		s.WsToSerial <- wsMsg
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
