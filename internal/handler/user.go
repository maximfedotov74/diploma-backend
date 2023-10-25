package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/utils"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
)

func (h *Handler) initUsersRoutes(router fiber.Router) {
	user := router.Group("/user")
	{
		user.Post("/registration", h.registration)
		user.Post("/login", h.login)
		user.Get("/:lk", h.authGuard, h.getLk)
		user.Get("/activate/:activationLink", h.activate)
		user.Get("/by-id/:id", h.getUserById)
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

	resp, appErr := h.services.UserService.Login(dto)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	access_cookie := new(fiber.Cookie)
	access_cookie.Name = "access_token"
	access_cookie.Value = resp.Tokens.AccessToken
	access_cookie.Expires = resp.Tokens.AccessExpTime
	refresh_cookie := new(fiber.Cookie)
	refresh_cookie.Name = "refresh_token"
	refresh_cookie.Value = resp.Tokens.RefreshToken
	refresh_cookie.Expires = resp.Tokens.RefreshExpTime
	refresh_cookie.HTTPOnly = true

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)

	return ctx.Status(201).JSON(resp)

}

func (h *Handler) activate(ctx *fiber.Ctx) error {
	activationLink := ctx.Params("activationLink")
	activated, err := h.services.UserService.Activate(activationLink)

	if err != nil || !activated {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Redirect(fmt.Sprintf("%s/api/user/lk", h.cfg.AppLink), 302)
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
	id, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	user, err := h.services.UserService.GetLk(*id)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(200).JSON(user)

}
