package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/utils"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

func (h *Handler) initUsersRoutes(router fiber.Router) {
	user := router.Group("/user")
	{
		user.Post("/registration", h.registration)
		user.Post("/login", h.login)
		user.Get("/activate/:activationLink", h.activate)
		user.Get("/by-id/:id", h.getUserById)
		user.Get("/refresh-token", h.refreshToken)
		user.Get("/:lk", h.authGuard, h.getLk)
	}
}

// @Summary Create user
// @Description Create user by body arguments
// @Tags users
// @Accept json
// @Produce json
// @Param dto body model.CreateUserDto true "create user with body dto"
// @Router /api/user/registration [post]
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

	id, appErr := h.services.UserService.Create(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(201).JSON(model.RegistrationResponse{Id: *id})
}

// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/user/by-id/:id [get]
// @Success 201 {object} model.User
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) getUserById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	user, appErr := h.services.UserService.GetUserById(id)

	if err != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(200).JSON(user)

}

// @Summary Login
// @Description Login to an account with account data
// @Tags users
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
	resp, appErr := h.services.UserService.Login(dto, userAgent)

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
// @Tags users
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
		appErr := lib.NewErr(messages.TOKEN_INVALID, 400)
		return ctx.Status(401).JSON(appErr)
	}

	userAgent := ctx.Get("User-Agent")

	response, err := h.services.UserService.RefreshToken(refreshToken, userAgent)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	access_cookie, refresh_cookie := utils.SetCookies(response.Tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)
	return ctx.Status(202).JSON(response)
}

func (h *Handler) activate(ctx *fiber.Ctx) error {
	activationLink := ctx.Params("activationLink")
	err := h.services.UserService.Activate(activationLink)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Redirect(fmt.Sprintf("https://ya.ru/"), 302)
}

// @Summary Get profile info
// @Security BearerToken
// @Description Get profile info by auth only
// @Tags users
// @Accept json
// @Produce json
// @Router /api/user/lk [get]
// @Success 200 {object} model.User
// @Failure 404 {object} lib.AppErr
// @Failure 401 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) getLk(ctx *fiber.Ctx) error {
	claims, err := utils.GetUserDataFromCtx(ctx)
	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	user, err := h.services.UserService.GetLk(claims.UserId)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(200).JSON(user)

}
