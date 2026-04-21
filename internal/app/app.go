package app

import (
	"github.com/bert1727/ChatApp/internal/handlers"
	"github.com/bert1727/ChatApp/internal/middlewares"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	"github.com/gofiber/fiber/v3"
)

func Build(app *fiber.App) {
	middlewares.SetupCustomMiddlewares(app)

	// Services
	userService := service.NewUserService()
	messageService := service.NewMessageService()

	hub := ws.NewHub(userService, messageService)
	authService := service.NewAuthService()

	go hub.Run()

	handlers := handlers.NewHandlers(authService, hub, userService)
	handlers.SetupHandlers(app)
	handlers.SetupRestrictedHandlers(app)
}
