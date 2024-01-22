package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type LocalSession struct {
	UserId    int
	UserAgent string
	Roles     []string
}

const LocalSessionKey = "user_session"

func GetLocalSession(ctx *fiber.Ctx) (*LocalSession, fall.Error) {
	data := ctx.Locals(LocalSessionKey)

	claims, ok := data.(LocalSession)
	if !ok {
		return nil, fall.NewErr(fall.UNAUTHORIZED, fall.STATUS_UNAUTHORIZED)
	}

	return &claims, nil
}
