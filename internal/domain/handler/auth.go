package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type authService interface {
	Login(ctx context.Context, dto model.LoginDto, userAgent string) (*model.LoginResponse, fall.Error)
	Registration(ctx context.Context, dto model.CreateUserDto) (*int, fall.Error)
	Refresh(ctx context.Context, refreshToken string, userAgent string) (*model.LoginResponse, fall.Error)
}

type AuthHandler struct {
	service        authService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewAuthHandler(service authService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *AuthHandler {
	return &AuthHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
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
// @Param dto body model.CreateUserDto true "Registation user with body dto"
// @Router /api/auth/registration [post]
// @Success 201 {object} model.RegistrationResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (ah *AuthHandler) registration(ctx *fiber.Ctx) error {
	dto := model.CreateUserDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {

		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	id, appErr := ah.service.Registration(ctx.Context(), dto)

	if appErr != nil {

		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(fall.STATUS_CREATED).JSON(model.RegistrationResponse{Id: *id})
}

// @Summary Login
// @Description Login to an account with account data
// @Tags auth
// @Accept json
// @Produce json
// @Param dto body model.LoginDto true "login in account"
// @Router /api/auth/login [post]
// @Success 201 {object} model.LoginResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (ah *AuthHandler) login(ctx *fiber.Ctx) error {

	var dto model.LoginDto

	err := ctx.BodyParser(&dto)
	if err != nil {
		ex := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	userAgent := ctx.Get("User-Agent")
	resp, appErr := ah.service.Login(ctx.Context(), dto, userAgent)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	access_cookie, refresh_cookie := utils.SetCookies(resp.Tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)

	return ctx.Status(fall.STATUS_CREATED).JSON(resp)

}

// @Summary Refresh tokens
// @Description Refresh tokens by cookies refresh_token
// @Tags auth
// @Accept json
// @Produce json
// @Router /api/auth/refresh-token [get]
// @Success 200 {object} model.LoginResponse
// @Failure 404 {object} fall.AppErr
// @Failure 401 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (ah *AuthHandler) refreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		appErr := fall.NewErr(fall.UNAUTHORIZED, fall.STATUS_UNAUTHORIZED)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	userAgent := ctx.Get("User-Agent")

	response, err := ah.service.Refresh(ctx.Context(), refreshToken, userAgent)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	access_cookie, refresh_cookie := utils.SetCookies(response.Tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)
	return ctx.Status(fall.STATUS_ACCEPTED).JSON(response)
}
