package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

func (h *Handler) authGuard(ctx *fiber.Ctx) error {
	ctx.Locals("userId", nil)

	authHeader := ctx.Get("Authorization")

	splittedHeader := strings.Split(authHeader, " ")

	if len(splittedHeader) != 2 {
		return ctx.Status(401).SendString("Unauthorized!")
	}

	accessToken := splittedHeader[1]

	userId, err := h.services.TokenService.Parse(accessToken, token.AccessToken)
	if err != nil {
		return ctx.Status(401).SendString("Unauthorized!")
	}

	ctx.Locals("userId", userId)
	return ctx.Next()

}
