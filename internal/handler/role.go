package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
)

// h.authGuard, h.roleGuard("ADMIN")
func (h *Handler) initRoleRoutes(router fiber.Router) {
	role := router.Group("/role")
	{
		role.Post("/", h.createRole)
		role.Post("/add-to-user", h.addRoleToUser)
		role.Delete("/remove-from-user", h.removeRoleFromUser)
	}
}

// @Summary Create role
// @Security BearerToken
// @Description Create role by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body model.CreateRoleDto true "create role with body dto"
// @Router /api/role/ [post]
// @Success 201 {object} model.Role
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {array} lib.AppErr
// @Failure 500 {array} lib.AppErr
func (h *Handler) createRole(ctx *fiber.Ctx) error {

	body := model.CreateRoleDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := lib.NewErr(err.Error(), 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	role, appErr := h.services.RoleService.Create(body)

	if appErr != nil {
		return ctx.Status(appErr.Status()).SendString(appErr.Message())
	}

	return ctx.Status(201).JSON(role)
}

// @Summary Add role to user
// @Security BearerToken
// @Description Add role to user by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body model.AddRoleToUserDto true "add role to user with body dto"
// @Router /api/role/add-to-user [post]
// @Success 201 {string} string
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {array} lib.AppErr
// @Failure 500 {array} lib.AppErr
func (h *Handler) addRoleToUser(ctx *fiber.Ctx) error {
	body := model.AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := lib.NewErr(err.Error(), 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := lib.ValidationMessages(error_messages)
		validError := lib.NewValidErr(items)
		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.services.RoleService.AddRoleToUser(body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(201).JSON(model.СompletedOperation{Completed: true})

}

// @Summary Remove role from user
// @Security BearerToken
// @Description Remove role from user by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body model.AddRoleToUserDto true "Remove role from user with body dto"
// @Router /api/role/remove-from-user [delete]
// @Success 201 {string} string
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {array} lib.AppErr
// @Failure 500 {array} lib.AppErr
func (h *Handler) removeRoleFromUser(ctx *fiber.Ctx) error {
	body := model.AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := lib.NewErr(err.Error(), 400)
		return ctx.Status(appErr.Status()).JSON(appErr.Message())
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := lib.ValidationMessages(error_messages)
		validError := lib.NewValidErr(items)
		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.services.RoleService.RemoveRoleFromUser(body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).SendString(appErr.Message())
	}

	return ctx.Status(200).JSON(model.СompletedOperation{Completed: true})

}
