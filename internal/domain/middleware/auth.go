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
	FindByAgentAndUserId(ctx context.Context, agent string, userId int) (*model.Session, fall.Error)
}

type userService interface {
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
}

func CreateAuthMiddleware(session sessionService, user userService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(utils.LocalSessionKey, nil)

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

		currentUser, _ := user.FindById(ctx.Context(), claims.UserId)
		if currentUser == nil {
			return ctx.Status(authError.Status()).JSON(authError)
		}

		session, ex := session.FindByAgentAndUserId(ctx.Context(), claims.UserAgent, claims.UserId)

		if ex != nil {
			access, refresh := utils.RemoveCookies()
			ctx.Cookie(access)
			ctx.Cookie(refresh)
			forbidded := fall.NewErr(fall.FORBIDDEN, fall.STATUS_FORBIDDEN)
			return ctx.Status(forbidded.Status()).JSON(forbidded)
		}

		contextData := model.LocalSession{UserId: session.UserId, UserAgent: session.UserAgent,
			Roles: currentUser.Roles, Email: currentUser.Email}

		ctx.Locals(utils.LocalSessionKey, contextData)
		return ctx.Next()
	}
}
