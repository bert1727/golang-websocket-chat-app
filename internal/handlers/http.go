package handlers

import (
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type HTTPHandler interface {
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
	RefreshToken(c fiber.Ctx) error
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
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *httpHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	log.Info().Str("refresh_token", refreshToken).Msg("refresh_token from cookie is")
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "there are no cookie with refresh token",
		})
	}

	newAccessToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": "refresh token is invalid",
		})
	}

	return c.JSON(fiber.Map{
		"newAccessToken": newAccessToken,
	})
}
