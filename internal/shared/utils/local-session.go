package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

const LocalSessionKey = "user_session"

func GetLocalSession(ctx *fiber.Ctx) (*model.LocalSession, fall.Error) {
	data := ctx.Locals(LocalSessionKey)

	claims, ok := data.(model.LocalSession)
	if !ok {
		return nil, fall.NewErr(fall.UNAUTHORIZED, fall.STATUS_UNAUTHORIZED)
	}

	return &claims, nil
}
