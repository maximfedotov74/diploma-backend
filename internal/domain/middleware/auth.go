package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
	"github.com/maximfedotov74/diploma-backend/internal/shared/keys"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type sessionService interface {
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, fall.Error)
}

type userService interface {
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
}

func CreateAuthMiddleware(session sessionService, user userService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(utils.LocalSessionKey, nil)
		userAgent := ctx.Get(keys.UserAgentHeader)

		authHeader := ctx.Get(keys.AuthorizationHeader)

		splittedHeader := strings.Split(authHeader, " ")

		authError := fall.NewErr(fall.UNAUTHORIZED, fall.STATUS_UNAUTHORIZED)

		if len(splittedHeader) != 2 {
			return ctx.Status(authError.Status()).JSON(authError)
		}

		accessToken := splittedHeader[1]
		claims, err := session.Parse(accessToken, jwt.AccessToken)

		if err != nil {
			return ctx.Status(authError.Status()).JSON(authError)
		}

		if userAgent != claims.UserAgent {
			return ctx.Status(authError.Status()).JSON(authError)
		}

		currentUser, _ := user.FindById(ctx.Context(), claims.UserId)
		if currentUser == nil {
			return ctx.Status(authError.Status()).JSON(authError)
		}

		contextData := model.LocalSession{UserId: currentUser.Id, UserAgent: claims.UserAgent, Roles: currentUser.Roles}

		ctx.Locals(utils.LocalSessionKey, contextData)
		return ctx.Next()
	}
}
