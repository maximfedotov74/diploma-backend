package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/model"
)

func (h *Handler) initRoleRoutes(router fiber.Router) {
	role := router.Group("/role")
	{
		role.Post("/", h.createRole)
		role.Post("/add-to-user", h.addRoleToUser)
		role.Delete("/remove-from-user", h.removeRoleFromUser)
	}
}

func (h *Handler) createRole(ctx *fiber.Ctx) error {

	body := model.CreateRoleDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	role, err := h.services.RoleService.Create(body)

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	return ctx.Status(201).JSON(role)
}

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

	flag, err := h.services.RoleService.AddRoleToUser(body.Title, body.UserId)

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	return ctx.Status(201).JSON(flag)

}

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

	flag, err := h.services.RoleService.RemoveRoleFromUser(body.Title, body.UserId)

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	return ctx.Status(200).JSON(flag)

}
