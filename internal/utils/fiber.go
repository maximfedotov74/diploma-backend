package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/constants"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

func GetUserDataFromCtx(ctx *fiber.Ctx) (*model.UserContextData, lib.Error) {
	data := ctx.Locals(constants.USER_CTX_KEY)

	claims, ok := data.(model.UserContextData)
	if !ok {
		return nil, lib.NewErr(messages.UNAUTHORIZED, 401)
	}

	return &claims, nil
}
