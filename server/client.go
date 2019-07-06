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
	if name == "" {
		panic("newClient: name is empty")
	}
	if hub == nil {
		panic("newClient: hub is nil")
	}
	if conn == nil {
		panic("newClient: conn is nil")
	}
	return &client{name, hub, conn}
}

func (c *client) send(message []byte) error {
	if &c == nil {
		panic("client.send(): client is nil")
	}

	return c.conn.WriteMessage(websocket.TextMessage, message)
}

func (c *client) listen() {
	if &c == nil {
		panic("client.listen(): client is nil")
	}

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
	if &c == nil {
		panic("client.close(): client is nil")
	}

	return c.conn.Close()
}
