package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type Client struct {
	Conn *websocket.Conn
	Send chan models.WSMessage
	done chan struct{}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan models.WSMessage),
		done: make(chan struct{}),
	}
}

func (c *Client) ListenRead(handler func([]byte)) {
	defer close(c.done)
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			return
		}
		handler(msg)
	}
}

func (c *Client) ListenWrite() {
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("Can not marshal data:", err)
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("WS write error:", err)
				return
			}
		case <-c.done:
			return
		}
	}
}
