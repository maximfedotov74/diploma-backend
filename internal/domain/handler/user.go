package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type userService interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
	Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error
	CreateChangePasswordCode(ctx context.Context, userId int) fall.Error
	ConfirmChangePassword(ctx context.Context, code string, userId int) fall.Error
	ChangePassword(ctx context.Context, dto model.ChangePasswordDto, localSession model.LocalSession) (*jwt.Tokens, fall.Error)
	GetUserSessions(ctx context.Context, userId int, agent string) (*model.UserSessionsResponse, fall.Error)
	GetAll(ctx context.Context, page int) (*model.GetAllUsersResponse, fall.Error)

	RemoveAllSessions(ctx context.Context, userId int) fall.Error
	RemoveSession(ctx context.Context, userId int, sessionId int) fall.Error
	RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) fall.Error
}

type UserHandler struct {
	service        userService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
	clientAuthUrl  string
}

func NewUserHandler(service userService, router fiber.Router, authMiddleware middleware.AuthMiddleware, clientAuthUrl string,

) *UserHandler {
	return &UserHandler{service: service, router: router, authMiddleware: authMiddleware, clientAuthUrl: clientAuthUrl}
}

func (h *UserHandler) InitRoutes() {
	userRouter := h.router.Group("/user")
	{
		userRouter.Get("/all", h.getAll)
		userRouter.Get("/session", h.authMiddleware, h.getSession)
		userRouter.Get("/session/all", h.authMiddleware, h.getSessions)
		userRouter.Get("/profile", h.authMiddleware, h.getProfile)
		userRouter.Patch("/profile", h.authMiddleware, h.updateProfile)
		userRouter.Patch("/password-code/change", h.authMiddleware, h.changePassword)
		userRouter.Post("/password-code/confirm", h.authMiddleware, h.confirmChangePasswordCode)
		userRouter.Post("/password-code", h.authMiddleware, h.createChangePasswordCode)

		userRouter.Delete("/session/all", h.authMiddleware, h.removeAllSessions)
		userRouter.Delete("/session/except-current/:sessionId", h.authMiddleware, h.removeExceptCurrentSession)
		userRouter.Delete("/session/:sessionId", h.authMiddleware, h.removeSession)
	}
}

// @Summary Get all users
// @Security BearerToken
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/all [get]
// @Param page query int false "Page"
// @Success 200 {object} model.GetAllUsersResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) getAll(ctx *fiber.Ctx) error {

	page := ctx.QueryInt("page", 1)

	users, err := h.service.GetAll(ctx.Context(), page)
	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(users)
}

// @Summary Remove all user sessions
// @Security BearerToken
// @Description Remove all user sessions
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/session/all [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) removeAllSessions(ctx *fiber.Ctx) error {

	session, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	ex = h.service.RemoveAllSessions(ctx.Context(), session.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	access_cookie, refresh_cookie := utils.RemoveCookies()

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Remove user sessions except current
// @Security BearerToken
// @Description Remove user sessions except current
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path int true "User session id"
// @Router /api/user/session/except-current/{sessionId} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) removeExceptCurrentSession(ctx *fiber.Ctx) error {

	session, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	sessionId, err := ctx.ParamsInt("sessionId")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex = h.service.RemoveExceptCurrentSession(ctx.Context(), session.UserId, sessionId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Remove user session by session id
// @Security BearerToken
// @Description Remove user session by session id
// @Tags user
// @Accept json
// @Produce json
// @Param sessionId path int true "User session id"
// @Router /api/user/session/{sessionId} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) removeSession(ctx *fiber.Ctx) error {
	session, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	sessionId, err := ctx.ParamsInt("sessionId")

	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex = h.service.RemoveSession(ctx.Context(), session.UserId, sessionId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Get base profile info
// @Security BearerToken
// @Description Get base profile info
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/profile [get]
// @Success 200 {object} model.User
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) getProfile(ctx *fiber.Ctx) error {

	session, err := utils.GetLocalSession(ctx)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	user, err := h.service.FindById(ctx.Context(), session.UserId)
	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(user)
}

// @Summary Get user sessions
// @Security BearerToken
// @Description Get user sessions
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/session/all [get]
// @Success 200 {object} model.UserSessionsResponse
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) getSessions(ctx *fiber.Ctx) error {

	session, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	refreshToken := ctx.Cookies("refresh_token")

	if refreshToken == "" {
		appErr := fall.NewErr(fall.UNAUTHORIZED, fall.STATUS_UNAUTHORIZED)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	sessions, ex := h.service.GetUserSessions(ctx.Context(), session.UserId, refreshToken)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(fall.STATUS_OK).JSON(sessions)
}

// @Summary Get local session
// @Security BearerToken
// @Description Get local session
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/session [get]
// @Success 200 {object} model.LocalSession
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) getSession(ctx *fiber.Ctx) error {

	session, err := utils.GetLocalSession(ctx)

	if err != nil {
		return ctx.Status(err.Status()).JSON(err)
	}

	return ctx.Status(fall.STATUS_OK).JSON(session)
}

// @Summary Update user profile info
// @Description Update user profile info
// @Tags user
// @Accept json
// @Produce json
// @Param dto body model.UpdateUserDto true "Update user profile with body dto"
// @Router /api/user/profile [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) updateProfile(ctx *fiber.Ctx) error {

	session, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.UpdateUserDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	validate.RegisterValidation("userGenderEnumValidation", model.UserGenderEnumValidation)

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	ex = h.service.Update(ctx.Context(), dto, session.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)

}

// @Summary Create change password code
// @Description Create change password code
// @Tags user
// @Accept json
// @Produce json
// @Router /api/user/password-code [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) createChangePasswordCode(ctx *fiber.Ctx) error {
	localSession, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	ex = h.service.CreateChangePasswordCode(ctx.Context(), localSession.UserId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	created := fall.GetCreated()
	return ctx.Status(created.Status()).JSON(created)
}

// @Summary Confirm change password code
// @Description Confirm change password code
// @Tags user
// @Accept json
// @Produce json
// @Param dto body model.ConfirmChangePasswordDto true "Confirm change password code with body dto"
// @Router /api/user/password-code/confirm [post]
// @Success 201 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) confirmChangePasswordCode(ctx *fiber.Ctx) error {
	localSession, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.ConfirmChangePasswordDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		ex := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	ex = h.service.ConfirmChangePassword(ctx.Context(), dto.Code, localSession.UserId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	created := fall.GetOk()
	return ctx.Status(created.Status()).JSON(created)
}

// @Summary Change password
// @Description Change password
// @Tags user
// @Accept json
// @Produce json
// @Param dto body model.ChangePasswordDto true "Change password with body dto"
// @Router /api/user/password-code/change [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *UserHandler) changePassword(ctx *fiber.Ctx) error {
	localSession, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.ChangePasswordDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		ex := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(fall.STATUS_BAD_REQUEST).JSON(validError)
	}

	tokens, ex := h.service.ChangePassword(ctx.Context(), dto, *localSession)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	access_cookie, refresh_cookie := utils.SetCookies(*tokens)

	ctx.Cookie(access_cookie)
	ctx.Cookie(refresh_cookie)

	ok := fall.GetOk()
	return ctx.Status(ok.Status()).JSON(ok)
}
