package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/utils"
)

func (h *Handler) initUsersRoutes(router fiber.Router) {
	user := router.Group("/user")
	{
		user.Get("/activate/:activationLink", h.activate)
		user.Get("/by-id/:id", h.getUserById)
		user.Get("/:lk", h.authGuard, h.getLk)
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

	if err != nil {
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

	user, err := h.services.UserService.GetUserById(claims.UserId)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(200).JSON(user)

}
