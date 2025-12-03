package ws

import (
	"fmt"

	"github.com/bert1727/ChatApp/internal/models"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/gofiber/fiber/v2/log"
)

type Hub interface {
	Run()
	SendToRoom(d string, msg models.Message)
	SendToUser(userID uint, msg models.Message)
	UnregisterClient(c *Client)
	RegisterClient(c *Client)
}

type hub struct {
	clients        map[uint]*Client
	register       chan *Client
	unregister     chan *Client
	broadcast      chan models.Message
	messageService service.MessageService
	userService    service.UserService
}

func NewHub(us service.UserService, ms service.MessageService) Hub {
	return &hub{
		clients:        make(map[uint]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan models.Message),
		userService:    us,
		messageService: ms,
	}
}

// func NewWSHandler(h *Hub) *WSHandler

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			fmt.Println("Client registered in hub, client id:", client.ID)

		case client := <-h.unregister: // unregister user when disconnected
			close(client.SendChan)
			delete(h.clients, client.ID)

		case msg := <-h.broadcast:
			if msg.ReceiverID != 0 {
				h.SendToUser(msg.ReceiverID, msg)
				continue
			}
			// if msg.RoomID != "" {
			// 	h.SendToRoom(msg.RoomID, msg)
			// }
		}
	}
}

func (h *hub) SendToRoom(d string, msg models.Message) {
	panic("unimplemented")
}

func (h *hub) SendToUser(userID uint, msg models.Message) {
	log.Info("the message was get almost... User id is ", userID)
	if c, ok := h.clients[userID]; ok {
		log.Info("the message sent to user:", "userID", userID, "msg", msg)
		c.Send(msg)
	}
}

func (h *hub) UnregisterClient(c *Client) {
	if _, ok := h.clients[c.ID]; ok {
		close(c.SendChan)
		delete(h.clients, c.ID)
	}
}

func (h *hub) RegisterClient(c *Client) {
	h.register <- c
}
