package main

import (
	"os"

	// "github.com/gofiber/contrib/v3/websocket"

	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/db"
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/handlers"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	"github.com/gofiber/contrib/v3/websocket"
	middleware "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	jwtware "github.com/gofiber/contrib/v3/jwt"
)

func main() {
	app := fiber.New()
	cfg := config.New()

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: "24",
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	log.Logger = zerolog.New(consoleWriter).With().Caller().Timestamp().Caller().Logger()
	app.Use(middleware.New(middleware.Config{
		Logger: &logger,
	}))

	if err := godotenv.Load(); err != nil {
		log.Info().Msg("cannot load .env file")
	}

	// logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	// app.Use(logger.New(logger.Config{
	// 	Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
	// }))

	app.Use(cors.New())

	// TODO: do this in separate function
	app.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		log.Info().Msg("User connected without ws")
		return fiber.NewError(fiber.StatusUpgradeRequired, "WebSocket upgrade required")
	})

	db := db.SetupDB()

	// Repositories
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)

	// Handlers
	hub := ws.NewHub(userService, messageService)
	wsHandler := handlers.NewWSHandler(hub, userService)

	authService := service.NewAuthService(userRepo, cfg)
	httpH := handlers.NewHTTPNandler(authService)
	authHandler := handlers.HTTPHandler(httpH)
	go hub.Run()

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	// Endpoints
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
		Extractor:  extractors.FromAuthHeader("Bearer"),
		Claims:     &domain.JWTClaims{},
	}))

	// Restricted
	app.Get("/ws", wsHandler.HandleWSConnection)

	// Test handler
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

	log.Info().Msg("Server running on http://localhost:6969")
	log.Fatal().Err(app.Listen(":6969"))
}
