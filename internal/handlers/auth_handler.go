package handlers

import (
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type AuthJWTHandler interface {
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
	RefreshToken(c fiber.Ctx) error
}

func NewAuthJWTHandler(s service.AuthService) AuthJWTHandler {
	return &authHandler{authService: s}
}

type authHandler struct {
	authService service.AuthService
}

func (h *authHandler) Register(c fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		log.Err(err).Msg("failed to parse data")
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError(
			"invalid data",
			fiber.StatusBadRequest,
			nil,
		))
	}

	res, err := h.authService.Register(&req)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to register user")

		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError(
			"user is not register",
			fiber.StatusBadRequest,
			nil,
		))
	}

	return c.JSON(res)
}

func (h *authHandler) Login(c fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError(
			"failed to parse request data",
			fiber.StatusBadRequest,
			nil,
		))
	}

	res, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError(
			"failed to login a user",
			fiber.StatusBadRequest,
			nil,
		))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HTTPOnly: true,
		Secure:   false,            // TODO: change for https
		SameSite: "Lax",            // NOTE: maybe change to "Lax"
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *authHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		log.Info().Msg("field refresh_token is empty in cookie ")
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewAPIError(
			"there are no cookie with refresh token",
			fiber.StatusBadRequest,
			nil,
		))
	}

	newAccessToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewAPIError(
			"refresh token is invalid",
			fiber.StatusUnauthorized,
			nil,
		))
	}

	return c.JSON(fiber.Map{
		"newAccessToken": newAccessToken,
	})
}
