package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/constants"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

func GetUserIdFromCtx(ctx *fiber.Ctx) (*int, lib.Error) {
	data := ctx.Locals(constants.USER_CTX_KEY)

	userId, ok := data.(int)

	if !ok {
		return nil, lib.NewErr(messages.UNAUTHORIZED, 401)
	}
	return &userId, nil
}
