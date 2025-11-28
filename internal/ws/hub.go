package ws

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

type Hub struct {
	clients        map[string]*Client
	Register       chan *Client
	Unregister     chan *Client
	broadcast      chan Message
	messageService *MessageService
}

type WSHandler struct {
	hub *Hub
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		// rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

// func NewWSHandler(h *Hub) *WSHandler

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.ID] = client
			fmt.Println("Client registered in hub:", client.ID)

		case client := <-h.Unregister:
			close(client.SendChan)
			delete(h.clients, client.ID)
			// h.unregisterClient(client)

		case msg := <-h.broadcast:
			if msg.ReceiverID != "" {
				h.SendToUser(msg.ReceiverID, msg)
				continue
			}
			if msg.RoomID != "" {
				h.sendToRoom(msg.RoomID, msg)
			}
		}
	}
}

func (h *Hub) sendToRoom(d string, msg Message) {
	panic("unimplemented")
}

func (h *Hub) SendToUser(userID string, msg Message) {
	log.Info("the message was get almost... User id is " + userID)
	if c, ok := h.clients[userID]; ok {
		log.Info("the message sent to user:", "userID", userID, "msg", msg)
		c.Send(msg)
	}
}

func (h *Hub) UnregisterClient(c *Client) {
	if _, ok := h.clients[c.ID]; ok {
		close(c.SendChan)
		delete(h.clients, c.ID)
	}
}

func (h *Hub) RegisterClient(c *Client) {
	h.Register <- c
}
