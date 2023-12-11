package wish

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type RoleGuard func(roles ...string) fiber.Handler
type AuthGuard fiber.Handler

type Service interface {
	AddToCart(dto AddToCartDto, userId int) exception.Error
	DeleteFromCart(userId int, modelSizeId int) exception.Error
	IncreaseNumber(userId int, modelSizeId int) exception.Error
	ReduceNumber(userId int, modelSizeId int) exception.Error
	RemoveSeveralItems(userId int, modelSizesIds []int) exception.Error
	GetUserCart(userId int) ([]CartItem, exception.Error)
	ToggleWish(modelId int, userId int) exception.Error
}

type WishHandler struct {
	service   Service
	router    fiber.Router
	authGuard AuthGuard
	roleGuard RoleGuard
}

func NewWishHandler(service Service, router fiber.Router, authGuard AuthGuard) *WishHandler {
	return &WishHandler{
		service:   service,
		router:    router,
		authGuard: authGuard,
	}
}

func (wh *WishHandler) InitRoutes() {
	wishRouter := wh.router.Group("wish")
	{
		wishRouter.Post("/", wh.authGuard, wh.toggleWish)
		wishRouter.Get("/cart", wh.authGuard, wh.getUserCart)
		wishRouter.Post("/cart", wh.authGuard, wh.addToCart)
		wishRouter.Delete("/cart/:modelSizeId", wh.authGuard, wh.deleteFromCart)
		wishRouter.Patch("/cart/increase/:modelSizeId", wh.authGuard, wh.increaseNumber)
		wishRouter.Patch("/cart/reduce/:modelSizeId", wh.authGuard, wh.reduceNumber)
	}
}

func (wh *WishHandler) toggleWish(ctx *fiber.Ctx) error {
	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := AddToWishDto{}

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

	ex = wh.service.ToggleWish(dto.ModelId, user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(exception.STATUS_OK)
}

func (wh *WishHandler) getUserCart(ctx *fiber.Ctx) error {
	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	items, ex := wh.service.GetUserCart(user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(exception.STATUS_OK).JSON(items)
}

func (wh *WishHandler) deleteFromCart(ctx *fiber.Ctx) error {
	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	modelSizeId, err := ctx.ParamsInt("modelSizeId")
	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	ex = wh.service.DeleteFromCart(user.UserId, modelSizeId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(exception.STATUS_OK)
}

func (wh *WishHandler) increaseNumber(ctx *fiber.Ctx) error {
	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	modelSizeId, err := ctx.ParamsInt("modelSizeId")
	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	ex = wh.service.IncreaseNumber(user.UserId, modelSizeId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(exception.STATUS_OK)
}

func (wh *WishHandler) reduceNumber(ctx *fiber.Ctx) error {
	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	modelSizeId, err := ctx.ParamsInt("modelSizeId")
	if err != nil {
		appErr := exception.NewErr(messages.VALIDATION_ID, exception.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	ex = wh.service.ReduceNumber(user.UserId, modelSizeId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.SendStatus(exception.STATUS_OK)
}

func (wh *WishHandler) addToCart(ctx *fiber.Ctx) error {

	user, ex := utils.GetUserDataFromCtx(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := AddToCartDto{}

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

	ex = wh.service.AddToCart(dto, user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	return ctx.SendStatus(exception.STATUS_CREATED)
}
