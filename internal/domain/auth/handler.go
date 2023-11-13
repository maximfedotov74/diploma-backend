package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/domain/user"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Service interface {
	Login(dto LoginDto, userAgent string) (*LoginResponse, exception.Error)
	Registration(dto user.CreateUserDto) (*int, exception.Error)
	Refresh(refreshToken string, userAgent string) (*LoginResponse, exception.Error)
}

type AuthHandler struct {
	service   Service
	router    fiber.Router
	authGuard fiber.Handler
}

func NewAuthHandler(service Service, router fiber.Router, authGuard fiber.Handler) *AuthHandler {
	return &AuthHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
	}
}

func (ah *AuthHandler) InitRoutes() {
	authRouter := ah.router.Group("/auth")
	{
		authRouter.Post("/registration", ah.registration)
		authRouter.Post("/login", ah.login)
		authRouter.Get("/refresh-token", ah.refreshToken)
	}
}

// @Summary Registation user
// @Description Registation by body arguments
// @Tags auth
// @Accept json
// @Produce json
// @Param dto body user.CreateUserDto true "Registation user with body dto"
// @Router /api/auth/registration [post]
// @Success 201 {object} auth.RegistrationResponse
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ah *AuthHandler) registration(ctx *fiber.Ctx) error {
	dto := user.CreateUserDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)

		return ctx.Status(400).JSON(validError)
	}

	id, appErr := ah.service.Registration(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(201).JSON(RegistrationResponse{Id: *id})
}

// @Summary Login
// @Description Login to an account with account data
// @Tags auth
// @Accept json
// @Produce json
// @Param dto body auth.LoginDto true "login in account"
// @Router /api/user/login [post]
// @Success 201 {object} auth.LoginResponse
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ah *AuthHandler) login(ctx *fiber.Ctx) error {

	var dto LoginDto

	err := ctx.BodyParser(&dto)
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	userAgent := ctx.Get("User-Agent")
	resp, appErr := ah.service.Login(dto, userAgent)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	access_cookie, refresh_cookie := utils.SetCookies(resp.Tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)

	return ctx.Status(201).JSON(resp)

}

// @Summary Refresh tokens
// @Description Refresh tokens by cookies refresh_token
// @Tags auth
// @Accept json
// @Produce json
// @Router /api/user/refresh-token [get]
// @Success 200 {object} auth.LoginResponse
// @Failure 404 {object} exception.AppErr
// @Failure 401 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (ah *AuthHandler) refreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		appErr := exception.NewErr(messages.UNAUTHORIZED, 400)
		return ctx.Status(401).JSON(appErr)
	}

	userAgent := ctx.Get("User-Agent")

	response, err := ah.service.Refresh(refreshToken, userAgent)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	access_cookie, refresh_cookie := utils.SetCookies(response.Tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)
	return ctx.Status(202).JSON(response)
}
