package handlers

import (
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/gofiber/fiber/v3"
)

type HTTPHandler interface {
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
}

func NewHTTPNandler(s service.AuthService) HTTPHandler {
	return &httpHandler{authService: s}
}

type httpHandler struct {
	authService service.AuthService
}

func (h *httpHandler) Register(c fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
	}

	res, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError("user is not register", fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(res)
}

// TODO: change return types
func (h *httpHandler) Login(c fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
	}

	res, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError(
			"failed to login a user",
			fiber.StatusBadRequest,
			err.Error(),
		))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HTTPOnly: true,
		Secure:   true, // TODO: change later
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	return c.Status(fiber.StatusOK).JSON(res)
}
