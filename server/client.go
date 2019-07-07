package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	name string
	hub  *hub
	conn *websocket.Conn
}

func newClient(name string, hub *hub, conn *websocket.Conn) *client {
	client := &client{name, hub, conn}
	hub.register <- client
	return client
}

func (c *client) send(message []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, message)
}

func (c *client) listen() {
	for {
		messageType, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println(err)
			c.hub.unregister <- c
			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		c.hub.broadcast <- clientMessage{c, message}
	}
}

func (c *client) close() error {
	return c.conn.Close()
}
