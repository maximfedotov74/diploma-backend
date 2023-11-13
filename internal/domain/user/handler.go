package user

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/models"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Service interface {
	GetUserById(id int) (*User, exception.Error)
	Activate(activationLink string) exception.Error
	ChangePassword(dto ChangePasswordDto, contextData *models.UserContextData) (*jwt.Tokens, exception.Error)
	CreateChangePasswordCode(userId int) exception.Error
}

type UserHandler struct {
	service   Service
	router    fiber.Router
	authGuard fiber.Handler
}

func NewUserHandler(service Service, router fiber.Router, authGuard fiber.Handler) *UserHandler {
	return &UserHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
	}
}

func (uh *UserHandler) InitRoutes() {
	userRouter := uh.router.Group("/user")
	{
		userRouter.Get("/activate/:activationLink", uh.activate)
		userRouter.Get("/by-id/:id", uh.getUserById)
		userRouter.Get("/lk", uh.authGuard, uh.getLk)

		userRouter.Post("/create-change-password-code", uh.authGuard, uh.createChangePasswordCode)

		userRouter.Patch("/change-password", uh.authGuard, uh.changePassword)
	}
}

// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/user/by-id/:id [get]
// @Success 201 {object} user.User
// @Failure 400 {object} exception.ValidationError
// @Failure 404 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *UserHandler) getUserById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	user, appErr := h.service.GetUserById(id)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.Status(200).JSON(user)

}

func (h *UserHandler) activate(ctx *fiber.Ctx) error {
	activationLink := ctx.Params("activationLink")
	err := h.service.Activate(activationLink)

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
// @Success 200 {object} user.User
// @Failure 404 {object} exception.AppErr
// @Failure 401 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *UserHandler) getLk(ctx *fiber.Ctx) error {
	claims, err := utils.GetUserDataFromCtx(ctx)
	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	user, err := h.service.GetUserById(claims.UserId)

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
// @Param dto body user.ChangePasswordDto true "change password"
// @Router /api/user/change-password [patch]
// @Success 200
// @Failure 404 {object} exception.AppErr
// @Failure 401 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *UserHandler) changePassword(ctx *fiber.Ctx) error {

	var dto ChangePasswordDto

	err := ctx.BodyParser(&dto)
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	claims, appErr := utils.GetUserDataFromCtx(ctx)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	tokens, appErr := h.service.ChangePassword(dto, claims)

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
// @Failure 404 {object} exception.AppErr
// @Failure 401 {object} exception.AppErr
// @Failure 500 {object} exception.AppErr
func (h *UserHandler) createChangePasswordCode(ctx *fiber.Ctx) error {

	claims, appErr := utils.GetUserDataFromCtx(ctx)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	appErr = h.service.CreateChangePasswordCode(claims.UserId)

	if appErr != nil {
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	return ctx.SendStatus(200)
}
