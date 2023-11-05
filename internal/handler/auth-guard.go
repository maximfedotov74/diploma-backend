package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/constants"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

func (h *Handler) authGuard(ctx *fiber.Ctx) error {
	ctx.Locals(constants.USER_CTX_KEY, nil)
	userAgent := ctx.Get("gent")

	authHeader := ctx.Get(constants.HEADER_AUTHORIZATION)

	splittedHeader := strings.Split(authHeader, " ")

	authError := lib.NewErr(messages.UNAUTHORIZED, 401)

	if len(splittedHeader) != 2 {
		return ctx.Status(authError.Status()).JSON(authError)
	}

	accessToken := splittedHeader[1]
	claims, err := h.services.TokenService.Parse(accessToken, token.AccessToken)
	if err != nil {
		return ctx.Status(authError.Status()).JSON(authError)
	}

	user, _ := h.services.UserService.GetUserById(claims.UserId)
	if user == nil {
		return ctx.Status(authError.Status()).JSON(authError)
	}

	contextData := model.UserContextData{User: *user, UserAgent: userAgent}

	ctx.Locals(constants.USER_CTX_KEY, contextData)
	return ctx.Next()

}
