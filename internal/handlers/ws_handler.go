package handlers

import (
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
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
	token := jwtware.FromContext(c)
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewAPIError(
			"missing token",
			fiber.StatusUnauthorized,
			nil,
		))
	}

	log.Info().Msg("token from context was successfully gotten")

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			domain.NewAPIError(
				"invalid token",
				fiber.StatusUnauthorized,
				nil,
			))
		// return c.SendStatus(fiber.StatusBadRequest).
	}

	user, err := h.userService.FindUserByID(claims.UserID)
	if err != nil {
		log.Info().Msg("user not found")

		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewAPIError(
			"unauthorized",
			fiber.StatusUnauthorized,
			nil,
		))
	}

	return websocket.New(func(conn *websocket.Conn) {
		client := &ws.Client{
			Conn:     conn,
			Hub:      h.hub,
			SendChan: make(chan domain.Message, 32),
			ID:       user.ID,
			Username: user.Username,
		}

		h.hub.RegisterClient(client)

		log.Info().
			Uint("client id", client.ID).
			Msg("client connected")

		defer func() {
			log.Info().
				Uint("client id", client.ID).
				Msg("client disconnected")

			h.hub.UnregisterClient(client)
			err := conn.Close()
			if err != nil {
				log.Err(err).Msg("can not close conn")
			}
		}()

		go client.WritePump()

		client.ReadPump()
	})(c)
}
