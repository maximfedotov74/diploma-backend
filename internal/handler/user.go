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
		user.Get("/activate/:activationLink", h.activate)
		user.Get("/by-id/:id", h.getUserById)
		user.Get("/lk", h.authGuard, h.getLk)

		user.Post("/create-change-password-code", h.authGuard, h.createChangePasswordCode)

		user.Patch("/change-password", h.authGuard, h.changePassword)

	}
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

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(200).JSON(user)

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

	user, err := h.services.UserService.GetUserById(claims.User.Id)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(200).JSON(user)

}

// @Summary Change password
// @Security BearerToken
// @Description Change user password
// @Tags users
// @Accept json
// @Produce json
// @Param dto body model.ChangePasswordDto true "change password"
// @Router /api/user/change-password [patch]
// @Success 200
// @Failure 404 {object} lib.AppErr
// @Failure 401 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) changePassword(ctx *fiber.Ctx) error {

	var dto model.ChangePasswordDto

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

	claims, appErr := utils.GetUserDataFromCtx(ctx)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	tokens, appErr := h.services.UserService.ChangePassword(dto, claims)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	access_cookie, refresh_cookie := utils.SetCookies(*tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)

	return ctx.SendStatus(200)
}

// @Summary Create change password code
// @Security BearerToken
// @Description Create change password code and send to email
// @Tags users
// @Accept json
// @Produce json
// @Router /api/user/create-change-password-code [post]
// @Success 200
// @Failure 404 {object} lib.AppErr
// @Failure 401 {object} lib.AppErr
// @Failure 500 {object} lib.AppErr
func (h *Handler) createChangePasswordCode(ctx *fiber.Ctx) error {

	claims, appErr := utils.GetUserDataFromCtx(ctx)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	appErr = h.services.UserService.CreateChangePasswordCode(claims.User)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(200)
}
