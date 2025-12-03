package handlers

import (
	"log"

	"github.com/bert1727/ChatApp/internal/models"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type WSHandler interface {
	HandleWSConnection(c fiber.Ctx) error
}

type wsHandler struct {
	hub         ws.Hub
	userService service.UserService
}

func NewWSHandler(h ws.Hub, us service.UserService) WSHandler {
	return &wsHandler{
		hub:         h,
		userService: us,
	}
}

func (h *wsHandler) HandleWSConnection(c fiber.Ctx) error {
	name := c.Query("name")
	password := c.Query("password")
	log.Println("params:", name, password)
	if name == "" {
		return c.Status(400).SendString("missing or invalid name")
	}

	user, err := h.userService.FindUserByName(name)
	log.Println("user was found")
	if err != nil {
		log.Println("user not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if user.Password != password {
		log.Println("missing or invalid password")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return websocket.New(func(conn *websocket.Conn) {
		client := &ws.Client{
			Conn:     conn,
			Hub:      h.hub,
			SendChan: make(chan models.Message, 32),
			ID:       user.ID,
		}

		h.hub.RegisterClient(client)
		log.Println("client connected:", client.ID)

		defer func() {
			log.Println("client disconnected:", client.ID)
			h.hub.UnregisterClient(client)
			conn.Close()
		}()

		go client.WritePump()
		client.ReadPump()
	})(c)
}

// TODO: complete this function

// func (h *wsHandler) HandleWSConnection(conn *websocket.Conn) {
// 	// user, err := h.userService.FindUserByID(userID)
// 	if err != nil {
// 		log.Println("invalid user")
// 		conn.Close()
//
// 		return
// 	}
//
// 	// client := &ws.Client{
// 	// 	Conn:     conn,
// 	// 	Hub:      h.hub,
// 	// 	SendChan: make(chan models.Message, 32), // limit the size
// 	// }
//
// 	log.Println("client connected")
// 	h.hub.RegisterClient(client)
//
// 	defer func() {
// 		log.Println("client disconnected:", client.ID)
// 		// client.Hub.Unregister <- client
// 		h.hub.UnregisterClient(client)
// 		conn.Close()
// 	}()
//
// 	go client.WritePump()
// 	client.ReadPump()
// }
