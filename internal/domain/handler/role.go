package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type roleService interface {
	AddRoleToUser(ctx context.Context, title string, userId int) fall.Error
	Create(ctx context.Context, dto model.CreateRoleDto) (*model.Role, fall.Error)
	RemoveRoleFromUser(ctx context.Context, title string, userId int) fall.Error
	FindWithUsers(ctx context.Context) ([]model.Role, fall.Error)
	FindRoleByTitle(ctx context.Context, title string) (*model.Role, fall.Error)
	RemoveRole(ctx context.Context, roleId int) fall.Error
	GetAll(ctx context.Context) ([]model.UserRole, fall.Error)
}

type RoleHandler struct {
	service        roleService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
	roleMiddleware middleware.RoleMiddleware
}

func NewRoleHandler(service roleService, router fiber.Router, authMiddleware middleware.AuthMiddleware,
	roleMiddleware middleware.RoleMiddleware) *RoleHandler {

	return &RoleHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
		roleMiddleware: roleMiddleware,
	}
}

func (rh *RoleHandler) InitRoutes() {
	roleRouter := rh.router.Group("/role")
	{
		roleRouter.Get("/", rh.getAll)
		roleRouter.Get("/with-users", rh.findAllWithRelations)
		roleRouter.Get("/:title", rh.findByTitle)
		roleRouter.Post("/", rh.createRole)
		roleRouter.Post("/add-to-user", rh.addRoleToUser)
		roleRouter.Delete("/remove-from-user", rh.removeRoleFromUser)
		roleRouter.Delete("/:id", rh.removeRole)
	}
}

// @Summary Remove role
// @Security BearerToken
// @Description Remove role by id
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "role id"
// @Router /api/role/{id} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) removeRole(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		validErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(validErr.Status()).JSON(validErr)
	}

	removeErr := h.service.RemoveRole(ctx.Context(), id)

	if removeErr != nil {
		return ctx.Status(removeErr.Status()).JSON(removeErr)
	}

	ok := fall.GetOk()

	return ctx.Status(ok.Status()).JSON(ok)

}

// @Summary Find all roles
// @Security BearerToken
// @Description Find all roles
// @Tags roles
// @Accept json
// @Produce json
// @Router /api/role/ [get]
// @Success 200 {array} model.UserRole
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) getAll(ctx *fiber.Ctx) error {

	roles, err := h.service.GetAll(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(roles)
}

// @Summary Find all roles with users
// @Security BearerToken
// @Description Find all roles with users
// @Tags roles
// @Accept json
// @Produce json
// @Router /api/role/with-uesers [get]
// @Success 200 {array} model.Role
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) findAllWithRelations(ctx *fiber.Ctx) error {

	roles, err := h.service.FindWithUsers(ctx.Context())

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(roles)
}

// @Summary Find one role by title
// @Description Find one role by title
// @Tags roles
// @Accept json
// @Produce json
// @Param title path string true "unique role title"
// @Router /api/role/{title} [get]
// @Success 200 {object} model.Role
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) findByTitle(ctx *fiber.Ctx) error {

	title := ctx.Params("title")

	role, err := h.service.FindRoleByTitle(ctx.Context(), title)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_CREATED).JSON(role)
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
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) createRole(ctx *fiber.Ctx) error {

	body := model.CreateRoleDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := fall.NewErr(err.Error(), fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	role, appErr := h.service.Create(ctx.Context(), body)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(fall.STATUS_CREATED).JSON(role)
}

// @Summary Add role to user
// @Security BearerToken
// @Description Add role to user by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body model.AddRoleToUserDto true "add role to user with body dto"
// @Router /api/role/add-to-user [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) addRoleToUser(ctx *fiber.Ctx) error {
	body := model.AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := fall.NewErr(err.Error(), fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)
		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.service.AddRoleToUser(ctx.Context(), body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	created := fall.GetCreated()

	return ctx.Status(created.Status()).JSON(created)
}

// @Summary Remove role from user
// @Security BearerToken
// @Description Remove role from user by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body model.AddRoleToUserDto true "Remove role from user with body dto"
// @Router /api/role/remove-from-user [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *RoleHandler) removeRoleFromUser(ctx *fiber.Ctx) error {
	body := model.AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := fall.NewErr(err.Error(), fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)
		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.service.RemoveRoleFromUser(ctx.Context(), body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ok := fall.GetOk()

	return ctx.Status(ok.Status()).JSON(ok)

}
