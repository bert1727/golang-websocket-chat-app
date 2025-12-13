/* Package ws: Represents a single WS connection */
package ws

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/gofiber/contrib/v3/websocket"
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
			c.Hub.UnregisterClient(c)
			slog.Error("Error reading json:", err)
			return
		}

		// forward to hub
		fmt.Println("Received message from client:", c.ID, "Message:", msg, msg.ReceiverID)
		c.Hub.SendToUser(msg.ReceiverID, msg)
		// c.Hub.broadcast <- msg
	}
}

func (c *Client) WritePump() {
	log.Println("write pump started for client:", c.ID)
	for msg := range c.SendChan {
		if err := c.Conn.WriteJSON(msg); err != nil {
			log.Println("Error writing json:", err)
			break
		}
	}
}

// func (c *Client) WritePump() {
// 	log.Println("Starting write pump for client:", c.ID)
// 	for {
// 		select {
// 		case msg, ok := <-c.SendChan:
// 			if !ok {
// 				// channel closed - client disconnected
// 				c.Conn.Close()
// 				log.Println("channel was closed")
// 				return
// 			}
// 			err := c.Conn.WriteJSON(msg)
// 			if err != nil {
// 				c.Conn.Close()
// 				return
// 			}
// 		}
// 	}
// }
//
// func NewClient(h *Hub, conn *websocket.Conn, userID uint) *Client
// func (c *Client) Send(msg []byte)
// func (c *Client) readPump()
// func (c *Client) writePump()
