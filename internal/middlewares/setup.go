package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func SetupCustomMiddlewares(app *fiber.App) {
	SetupCors(app)
	SetupRecover(app)
	SetWSConnectionMiddleware(app)
}
