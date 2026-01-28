package app

import (
	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/db"
	"github.com/bert1727/ChatApp/internal/handlers"
	"github.com/bert1727/ChatApp/internal/middlewares"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	"github.com/gofiber/fiber/v3"
)

func Build(app *fiber.App, cfg *config.Config) {
	middlewares.SetupCustomMiddlewares(app)

	db := db.SetupDB()

	// Repositories
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)
	hub := ws.NewHub(userService, messageService)
	authService := service.NewAuthService(userRepo, cfg)

	go hub.Run()

	handlers := handlers.NewHandlers(authService, hub, userService, cfg)
	handlers.SetupHandlers(app)
	handlers.SetupRestrictedHandlers(app)
}
