package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/api"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/opencv"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/rs/zerolog"
)

type Server struct {
	Clients map[*Client]bool
	dClient *device.DeviceClient
	mu      sync.RWMutex
	log     zerolog.Logger

	Send chan models.WSResponse

	repo   *repository.Repository
	infra  *infrastructure.Infrastructure
	camera *camera.Camera
	opencv *opencv.Classification
	api    *api.API

	router *WSRouter
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewServer(
	dClient *device.DeviceClient,
	repo *repository.Repository,
	url string,
	log zerolog.Logger,
	infra *infrastructure.Infrastructure,
	camera *camera.Camera,
	opencv *opencv.Classification,
) (*Server, error) {
	api := api.NewAPI(url)
	router := NewWsRouter()

	server := &Server{
		Clients: make(map[*Client]bool),
		dClient: dClient,
		log:     log,

		Send: make(chan models.WSResponse, 100),

		repo:   repo,
		infra:  infra,
		camera: camera,
		opencv: opencv,
		api:    api,

		router: router,
	}

	go server.ListenDeviceResponse()
	go server.ListenSend()

	return server, nil
}

func (s *Server) ListenDeviceResponse() {
	for data := range s.dClient.RespCh {
		s.broadcast(data)
	}
}

func (s *Server) ListenSend() {
	for msg := range s.Send {
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(msg models.WSResponse) {
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
		fmt.Println("WS INCOMING RAW")
		var wsMsg models.WSRequest
		if err := jsoniter.Unmarshal(msg, &wsMsg); err != nil {
			log.Println("WS to COM error:", err)
			return
		}

		switch wsMsg.Type {
		case "STOP":
			client.Control <- wsMsg
		case "RETRY", "SKIP", "ABORT":
			client.Control <- wsMsg
		default:
			fmt.Println("CHANNEL FULL")
		}

		handler, err := s.router.WsRouter(wsMsg)
		if err != nil {
			log.Println("WS to COM error:", err)
			return
		}

		go handler(s, client, wsMsg)
	})

	// Отправка сообщений клиенту
	go client.ListenWrite()

	// Обработка отключения
	go func() {
		<-client.done
		conn.Close()
		s.mu.Lock()
		s.dClient.Stop(client.SessionId)
		delete(s.Clients, client)
		s.mu.Unlock()
		log.Println("Client disconnected, total clients:", len(s.Clients))
	}()
}

func (s *Server) Close() {
	s.mu.Lock()
	for client := range s.Clients {
		close(client.Send)
		delete(s.Clients, client)
		client.Conn.Close()
	}
	s.mu.Unlock()
	log.Println("All users deleted")

	log.Println("COM port closed")
}
