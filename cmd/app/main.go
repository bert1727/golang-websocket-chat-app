package main

import (
	"errors"

	"github.com/bert1727/ChatApp/internal/app"
	"github.com/bert1727/ChatApp/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func main() {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).SendString("something went wrong")
		},
	})

	// logger.SetupLoggerWithFileConsoleWriter(fiberApp, "logs.log")
	logger.SetupLoggerWithConsoleWriter(fiberApp)

	app.Build(fiberApp)

	log.Info().Msg("Server running on http://localhost:6969")
	log.Fatal().Err(fiberApp.Listen(":6969"))
}
