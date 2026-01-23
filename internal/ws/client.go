/* Package ws: Represents a single WS connection */
package ws

import (
	"fmt"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/rs/zerolog/log"
)

// type client interface{}

// TODO: make an interface for client
type Client struct {
	ID       uint
	Conn     *websocket.Conn // placeholder for actual ws conn
	Hub      Hub
	SendChan chan domain.Message
	RoomID   string
}

func (c *Client) Send(msg domain.Message) {
	c.SendChan <- msg
}

func (c *Client) ReadPump() {
	for {
		var msg domain.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			c.sendError()
			c.Hub.UnregisterClient(c)
			log.Err(err).Msg("Failed to read json message")
			return
		}

		// forward to hub
		fmt.Println("Received message from client:", c.ID, "Message:", msg, msg.ReceiverID)
		c.Hub.SendToUser(msg.ReceiverID, msg)
		// c.Hub.broadcast <- msg
	}
}

func (c *Client) WritePump() {
	log.Info().Uint("client Id", c.ID).Msg("write pump started for client:")
	for msg := range c.SendChan {
		if err := c.Conn.WriteJSON(msg); err != nil {
			log.Err(err).Msg("Error writing json")
			break
		}
	}
}

func (c *Client) sendError() {
	errorMsg := domain.Message{
		SenderID:   0,
		ReceiverID: c.ID,
		Type:       "error",
		Content:    "failed to parse a message",
	}
	err := c.Conn.WriteJSON(errorMsg)
	if err != nil {
		log.Error().Msg("failed to write json error")
	}

	c.SendChan <- errorMsg
}
