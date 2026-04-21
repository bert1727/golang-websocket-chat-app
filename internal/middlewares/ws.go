package middlewares

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func SetupWSConnectionMiddleware(app *fiber.App) {
	app.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		log.Info().Msg("User connected without ws")
		return fiber.NewError(fiber.StatusUpgradeRequired, "WebSocket upgrade required")
	})
}
