package main

import (
	"fmt"
	"log"
)

type clientMessage struct {
	client  *client
	message []byte
}

type hub struct {
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
	broadcast  chan clientMessage
}

func newHub() *hub {
	return &hub{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan clientMessage),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("client %s registered\n", client.name)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Printf("client %s unregistered\n", client.name)
				err := client.close()
				if err != nil {
					log.Println(err)
				}
			}
		case clientMessage := <-h.broadcast:
			message := fmt.Sprintf("%s says: %s\n", clientMessage.client.name, string(clientMessage.message))
			for client := range h.clients {
				err := client.send([]byte(message))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
