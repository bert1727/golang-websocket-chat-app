package httphandlers

import (
	"fmt"
	"log"

	"github.com/bert1727/ChatApp/internal/service"
	"github.com/gofiber/fiber/v3"
)

type UserHandler interface {
	Register(c fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

func NewUserHanlder(us service.UserService) UserHandler {
	return &userHandler{userService: us}
}

type registerRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *userHandler) Register(c fiber.Ctx) error {
	var req registerRequest

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	user, err := h.userService.Register(req.Name, req.Password)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Println("user was added")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fmt.Sprintf("user was added with name: %s", user.Name),
	})
	// return c.JSON(user) // TODO: return only user id or name
}
