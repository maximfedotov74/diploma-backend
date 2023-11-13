package role

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/models"
)

type Service interface {
	AddRoleToUser(title string, userId int) exception.Error
	Create(dto CreateRoleDto) (*Role, exception.Error)
	RemoveRoleFromUser(title string, userId int) exception.Error
}

type RoleGuard interface {
	CheckRoles(roles ...string) fiber.Handler
}

type RoleHandler struct {
	service   Service
	router    fiber.Router
	authGuard fiber.Handler
	roleGuard RoleGuard
}

func NewRoleHandler(service Service, authGuard fiber.Handler, roleGuard RoleGuard, router fiber.Router) *RoleHandler {

	return &RoleHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
		roleGuard: roleGuard,
	}
}

func (rh *RoleHandler) InitRoutes() {
	roleRouter := rh.router.Group("/role")
	{
		roleRouter.Post("/", rh.authGuard, rh.roleGuard.CheckRoles("ADMIN"), rh.createRole)
		roleRouter.Post("/add-to-user", rh.addRoleToUser)
		roleRouter.Delete("/remove-from-user", rh.removeRoleFromUser)
	}
}

// @Summary Create role
// @Security BearerToken
// @Description Create role by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body role.CreateRoleDto true "create role with body dto"
// @Router /api/role/ [post]
// @Success 201 {object} role.Role
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *RoleHandler) createRole(ctx *fiber.Ctx) error {

	body := CreateRoleDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := exception.NewErr(err.Error(), 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	role, appErr := h.service.Create(body)

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
// @Param dto body role.AddRoleToUserDto true "add role to user with body dto"
// @Router /api/role/add-to-user [post]
// @Success 201 {string} models.СompletedOperation
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *RoleHandler) addRoleToUser(ctx *fiber.Ctx) error {
	body := AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := exception.NewErr(err.Error(), 400)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)
		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.service.AddRoleToUser(body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(201).JSON(models.СompletedOperation{Completed: true})

}

// @Summary Remove role from user
// @Security BearerToken
// @Description Remove role from user by body arguments
// @Tags roles
// @Accept json
// @Produce json
// @Param dto body role.AddRoleToUserDto true "Remove role from user with body dto"
// @Router /api/role/remove-from-user [delete]
// @Success 201 {string} models.СompletedOperation
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *RoleHandler) removeRoleFromUser(ctx *fiber.Ctx) error {
	body := AddRoleToUserDto{}

	err := ctx.BodyParser(&body)

	if err != nil {
		appErr := exception.NewErr(err.Error(), 400)
		return ctx.Status(appErr.Status()).JSON(appErr.Message())
	}

	validate := validator.New()

	err = validate.Struct(&body)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)
		return ctx.Status(validError.Status).JSON(validError)
	}

	appErr := h.service.RemoveRoleFromUser(body.Title, body.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).SendString(appErr.Message())
	}

	return ctx.Status(200).JSON(models.СompletedOperation{Completed: true})

}
