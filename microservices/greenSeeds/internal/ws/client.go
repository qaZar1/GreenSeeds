package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	done chan struct{}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan []byte, 10),
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
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("WS write error:", err)
				return
			}
		case <-c.done:
			return
		}
	}
}
