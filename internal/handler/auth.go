package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/utils"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

func (h *Handler) initAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	{
		auth.Post("/registration", h.registration)
		auth.Post("/login", h.login)
		auth.Get("/refresh-token", h.refreshToken)
		auth.Get("/ip", func(c *fiber.Ctx) error {

			ip := "94.181.43.165"

			return c.JSON(ip)
		})
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
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) registration(ctx *fiber.Ctx) error {
	dto := model.CreateUserDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := lib.ValidationMessages(error_messages)
		validError := lib.NewValidErr(items)

		return ctx.Status(400).JSON(validError)
	}

	id, appErr := h.services.AuthService.Registration(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(201).JSON(model.RegistrationResponse{Id: *id})
}

// @Summary Login
// @Description Login to an account with account data
// @Tags auth
// @Accept json
// @Produce json
// @Param dto body model.LoginDto true "login in account"
// @Router /api/user/login [post]
// @Success 201 {object} model.LoginResponse
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) login(ctx *fiber.Ctx) error {

	var dto model.LoginDto

	err := ctx.BodyParser(&dto)
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := lib.ValidationMessages(error_messages)
		validError := lib.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	userAgent := ctx.Get("User-Agent")
	resp, appErr := h.services.AuthService.Login(dto, userAgent)

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
// @Success 200 {object} model.LoginResponse
// @Failure 404 {object} lib.AppErr
// @Failure 401 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) refreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		appErr := lib.NewErr(messages.UNAUTHORIZED, 400)
		return ctx.Status(401).JSON(appErr)
	}

	userAgent := ctx.Get("User-Agent")

	response, err := h.services.AuthService.Refresh(refreshToken, userAgent)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	access_cookie, refresh_cookie := utils.SetCookies(response.Tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)
	return ctx.Status(202).JSON(response)
}
