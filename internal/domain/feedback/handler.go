package feedback

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/shared/constants"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Service interface {
	AddFeedback(dto AddFeedbackDto, userId int) exception.Error
	GetModelFeedback(modelId int, order string) (*ModelFeedbackResponse, exception.Error)
	GetAll(order string) ([]Feedback, exception.Error)
	DeleteFeedback(feedbackId int) exception.Error
	ToggleHidden(feedbackId int) exception.Error
}

type RoleGuard func(roles ...string) fiber.Handler
type AuthGuard fiber.Handler

type FeedbackHandler struct {
	service   Service
	router    fiber.Router
	authGuard AuthGuard
	roleGuard RoleGuard
}

func NewFeedbackHandler(service Service, router fiber.Router, authGuard AuthGuard, roleGuard RoleGuard) *FeedbackHandler {
	return &FeedbackHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
		roleGuard: roleGuard,
	}
}

func (fh *FeedbackHandler) InitRoutes() {
	feedbackRouter := fh.router.Group("feedback")
	{
		feedbackRouter.Post("/", fh.authGuard, fh.roleGuard(constants.ADMIN_ROLE), fh.addFeedback)
		feedbackRouter.Delete("/:id", fh.authGuard, fh.roleGuard(constants.ADMIN_ROLE), fh.deleteFeedback)
		feedbackRouter.Patch("/:id", fh.authGuard, fh.roleGuard(constants.ADMIN_ROLE), fh.toggleHidden)
		feedbackRouter.Get("/", fh.authGuard, fh.roleGuard(constants.ADMIN_ROLE), fh.getAll)
		feedbackRouter.Get("/:modelId", fh.getModelFeedback)
	}
}

func (fh *FeedbackHandler) toggleHidden(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := fh.service.ToggleHidden(id)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

func (fh *FeedbackHandler) getAll(ctx *fiber.Ctx) error {
	order := ctx.Query("order", "ASC")

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	feedback, ex := fh.service.GetAll(order)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(feedback)
}

func (fh *FeedbackHandler) deleteFeedback(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	ex := fh.service.DeleteFeedback(id)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_OK)
}

func (fh *FeedbackHandler) addFeedback(ctx *fiber.Ctx) error {

	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := AddFeedbackDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := exception.NewErr(messages.INVALID_BODY, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := exception.ValidationMessages(error_messages)
		validError := exception.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	ex = fh.service.AddFeedback(dto, user.UserId)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}

func (fh *FeedbackHandler) getModelFeedback(ctx *fiber.Ctx) error {

	modelId, err := ctx.ParamsInt("modelId")
	order := ctx.Query("order", "ASC")

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	feedback, ex := fh.service.GetModelFeedback(modelId, order)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.Status(exception.STATUS_OK).JSON(feedback)
}
