package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func SetupRecover(app *fiber.App) {
	app.Use(recover.New())
}
