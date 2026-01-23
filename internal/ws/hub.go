package ws

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

type Hub interface {
	Run()
	SendToRoom(d string, msg domain.Message)
	SendToUser(userID uint, msg domain.Message)
	UnregisterClient(c *Client)
	RegisterClient(c *Client)
}

type hub struct {
	clients        map[uint]*Client
	register       chan *Client
	unregister     chan *Client
	broadcast      chan domain.Message
	messageService service.MessageService
	userService    service.UserService
	mu             sync.RWMutex
}

func NewHub(us service.UserService, ms service.MessageService) Hub {
	return &hub{
		clients:        make(map[uint]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan domain.Message),
		userService:    us,
		messageService: ms,
	}
}

// func NewWSHandler(h *Hub) *WSHandler

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			fmt.Println("Client registered in hub, client id:", client.ID)

		case client := <-h.unregister: // unregister user when disconnected
			close(client.SendChan)
			delete(h.clients, client.ID)

		case msg := <-h.broadcast:
			if msg.ReceiverID != 0 {
				h.SendToUser(msg.ReceiverID, msg)
				continue
			}
		}
	}
}

func (h *hub) SendToRoom(d string, msg domain.Message) {
	panic("unimplemented")
}

func (h *hub) SendToUser(userID uint, msg domain.Message) {
	log.Info("the message was get almost... User id is ", userID)
	if c, ok := h.clients[userID]; ok {
		log.Info("the message sent to user:", "userID", userID, "msg", msg.Content)

		msg.SenderID = c.ID
		err := validator.New().Struct(msg)
		if err != nil {
			// log.Error().Err(err).Msg("Failed to validate a new message")

			var validateErrs validator.ValidationErrors
			if errors.As(err, &validateErrs) {
				for _, e := range validateErrs {
					log.Info("error is:", e)
					// log.Info().Msg(e.Namespace())
				}
			}

			c.sendError()
			return
		}

		if err := h.messageService.SaveMessage(&msg); err != nil {
			log.Info("failed to save a message in db")
		}

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
