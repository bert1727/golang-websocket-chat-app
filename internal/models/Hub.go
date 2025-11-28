package models

import (
	"log"

	"github.com/gofiber/contrib/v3/websocket"
)

type Hub struct {
	Clients   map[int]*websocket.Conn
	Broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[int]*websocket.Conn),
		Broadcast: make(chan []byte),
	}
}

func (h *Hub) Run(id int) {
	for {
		message := <-h.Broadcast
		for _, client := range h.Clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				client.Close()
				log.Println("client with id:", id, " was deleted ")
				delete(h.Clients, id)
			}
		}
	}
}
