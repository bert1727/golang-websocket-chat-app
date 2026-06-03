package handlers

import (
	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/rs/zerolog/log"
)

type Handlers struct {
	Auth AuthJWTHandler
	WS   WSHandler
	Cfg  *config.Config
}

func (h *Handlers) SetupHandlers(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.RefreshToken)
}

func (h *Handlers) SetupRestrictedHandlers(app *fiber.App) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(h.Cfg.JWTSecret)},
		// Extractor:  extractors.FromAuthHeader("Bearer"),
		Extractor: extractors.FromQuery("token"),
		Claims:    &domain.JWTClaims{},
	}))
	app.Get("/ws", h.WS.HandleWSConnection)
}

func NewHandlers(
	authService service.AuthService,
	hub ws.Hub,
	userService service.UserService,
) *Handlers {
	return &Handlers{
		Auth: NewAuthJWTHandler(authService),
		WS:   NewWSHandler(hub, userService),
		Cfg:  config.New(),
	}
}

func testingHandler(app *fiber.App) {
	app.Get("/test", func(c fiber.Ctx) error {
		token := jwtware.FromContext(c)
		claims := token.Claims.(*domain.JWTClaims)
		if claims == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		log.Info().Any("token", token).Msg("token is")
		log.Info().Uint("id", claims.UserID).Str("email", claims.Email).Msg("user data from token")

		if token != nil {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.SendStatus(fiber.StatusBadRequest)
	})
}
