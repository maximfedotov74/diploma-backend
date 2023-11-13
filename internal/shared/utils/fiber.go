package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/shared/constants"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/models"
)

func GetUserDataFromCtx(ctx *fiber.Ctx) (*models.UserContextData, exception.Error) {
	data := ctx.Locals(constants.USER_CTX_KEY)

	claims, ok := data.(models.UserContextData)
	if !ok {
		return nil, exception.NewErr(messages.UNAUTHORIZED, 401)
	}

	return &claims, nil
}
