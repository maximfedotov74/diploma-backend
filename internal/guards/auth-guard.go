package guards

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/domain/user"
	"github.com/maximfedotov74/fiber-psql/internal/shared/constants"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/models"

	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

type SessionService interface {
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, error)
}

type UserService interface {
	GetUserById(id int) (*user.User, exception.Error)
}

type AuthGuard struct {
	sessionService SessionService
	userService    UserService
}

func NewAuthGuard(sessionService SessionService, userService UserService) *AuthGuard {
	return &AuthGuard{
		sessionService: sessionService,
		userService:    userService,
	}
}

func (ag *AuthGuard) CheckAuth(ctx *fiber.Ctx) error {
	ctx.Locals(constants.USER_CTX_KEY, nil)
	userAgent := ctx.Get("User-Agent")

	authHeader := ctx.Get(constants.HEADER_AUTHORIZATION)

	splittedHeader := strings.Split(authHeader, " ")

	authError := exception.NewErr(messages.UNAUTHORIZED, 401)

	if len(splittedHeader) != 2 {
		return ctx.Status(authError.Status()).JSON(authError)
	}

	accessToken := splittedHeader[1]
	claims, err := ag.sessionService.Parse(accessToken, jwt.AccessToken)
	if err != nil {
		return ctx.Status(authError.Status()).JSON(authError)
	}

	currentUser, _ := ag.userService.GetUserById(claims.UserId)
	if currentUser == nil {
		return ctx.Status(authError.Status()).JSON(authError)
	}

	contextData := models.UserContextData{UserId: currentUser.Id, UserAgent: userAgent}

	ctx.Locals(constants.USER_CTX_KEY, contextData)
	return ctx.Next()
}
