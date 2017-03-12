package main

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	//List of currently connected clients
	clients map[*websocket.Conn]bool

	//Connections are sending message on this channel
	broadcast chan []byte

	//Request to be added to list of connections
	register chan *websocket.Conn

	//Request to disconnect
	unregister chan *websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					delete(h.clients, client)
				}
			}
		}
	}
}
