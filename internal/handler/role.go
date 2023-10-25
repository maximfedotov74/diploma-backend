package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
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
		return ctx.Status(400).SendString(err.Error())
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
// @Success 201 {boolean} bool
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {array} lib.AppErr
// @Failure 500 {array} lib.AppErr
func (h *Handler) addRoleToUser(ctx *fiber.Ctx) error {
	body := model.AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		return ctx.Status(400).SendString(fmt.Sprintf("%s", error_messages))
	}

	flag, appErr := h.services.RoleService.AddRoleToUser(body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).SendString(appErr.Message())
	}

	return ctx.Status(201).JSON(flag)

}

// @Summary Remove role from user
// @Security BearerToken
// @Description Remove role from user by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body model.AddRoleToUserDto true "Remove role from user with body dto"
// @Router /api/role/remove-from-user [delete]
// @Success 201 {boolean} bool
// @Failure 400 {object} lib.ValidationError
// @Failure 404 {array} lib.AppErr
// @Failure 500 {array} lib.AppErr
func (h *Handler) removeRoleFromUser(ctx *fiber.Ctx) error {
	body := model.AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		return ctx.Status(400).SendString(fmt.Sprintf("%s", error_messages))
	}

	flag, appErr := h.services.RoleService.RemoveRoleFromUser(body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).SendString(appErr.Message())
	}

	return ctx.Status(200).JSON(flag)

}
