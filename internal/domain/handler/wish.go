package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/domain/middleware"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

type wishService interface {
	GetUserWish(ctx context.Context, userId int) ([]*model.CatalogProductModel, fall.Error)
	FindModelInUserCart(ctx context.Context, modelSizeId int, userId int) (*model.CartItemModel, fall.Error)
	AddToCart(ctx context.Context, dto model.AddToCartDto, userId int) fall.Error
	DeleteFromCart(ctx context.Context, userId int, modelSizeId int) fall.Error
	IncreaseNumber(ctx context.Context, userId int, modelSizeId int) fall.Error
	ReduceNumber(ctx context.Context, userId int, modelSizeId int) fall.Error
	RemoveSeveralItems(ctx context.Context, userId int, modelSizesIds []int) fall.Error
	GetUserCart(ctx context.Context, userId int) ([]model.CartItem, fall.Error)
	ToggleWish(ctx context.Context, modelId int, userId int) fall.Error
}

type WishHandler struct {
	service        wishService
	router         fiber.Router
	authMiddleware middleware.AuthMiddleware
}

func NewWishHandler(service wishService, router fiber.Router, authMiddleware middleware.AuthMiddleware) *WishHandler {
	return &WishHandler{
		service:        service,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func (wh *WishHandler) InitRoutes() {
	wishRouter := wh.router.Group("wish")
	{
		wishRouter.Post("/cart", wh.authMiddleware, wh.addToCart)
		wishRouter.Post("/", wh.authMiddleware, wh.toggleWish)

		wishRouter.Get("/cart", wh.authMiddleware, wh.getUserCart)
		wishRouter.Get("/", wh.authMiddleware, wh.getUserWish)

		wishRouter.Patch("/cart/increase/:modelSizeId", wh.authMiddleware, wh.increaseNumber)
		wishRouter.Patch("/cart/reduce/:modelSizeId", wh.authMiddleware, wh.reduceNumber)

		wishRouter.Delete("/cart/several/:ids", wh.authMiddleware, wh.removeSeveralItems)
		wishRouter.Delete("/cart/:modelSizeId", wh.authMiddleware, wh.deleteFromCart)

	}
}

// @Summary Toggle wishlist item
// @Security BearerToken
// @Description Toggle wishlist item
// @Tags wish
// @Accept json
// @Produce json
// @Param dto body model.AddToWishDto true "Toggle with dto"
// @Router /api/wish/ [post]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) toggleWish(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.AddToWishDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	ex = wh.service.ToggleWish(ctx.Context(), dto.ModelId, user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Get user cart items
// @Security BearerToken
// @Description Get user cart items
// @Tags wish
// @Accept json
// @Produce json
// @Router /api/wish/cart [get]
// @Success 200 {array} model.CartItem
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) getUserCart(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	items, ex := wh.service.GetUserCart(ctx.Context(), user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(fall.STATUS_OK).JSON(items)
}

// @Summary Get user wish items
// @Security BearerToken
// @Description Get wish cart items
// @Tags wish
// @Accept json
// @Produce json
// @Router /api/wish/ [get]
// @Success 200 {array} model.CatalogProductModel
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) getUserWish(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	items, ex := wh.service.GetUserWish(ctx.Context(), user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	return ctx.Status(fall.STATUS_OK).JSON(items)
}

// @Summary Delete item from cart
// @Security BearerToken
// @Description Delete item from cart
// @Tags wish
// @Accept json
// @Produce json
// @Param modelSizeId path int true "Model Size Id"
// @Router /api/wish/cart/{modelSizeId} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) deleteFromCart(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	modelSizeId, err := ctx.ParamsInt("modelSizeId")
	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	ex = wh.service.DeleteFromCart(ctx.Context(), user.UserId, modelSizeId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Increase cart item quantity
// @Security BearerToken
// @Description Increase cart item quantity
// @Tags wish
// @Accept json
// @Produce json
// @Param modelSizeId path int true "Model Size Id"
// @Router /api/wish/cart/increase/{modelSizeId} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) increaseNumber(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	modelSizeId, err := ctx.ParamsInt("modelSizeId")
	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	ex = wh.service.IncreaseNumber(ctx.Context(), user.UserId, modelSizeId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Reduce cart item quantity
// @Security BearerToken
// @Description Increase cart item quantity
// @Tags wish
// @Accept json
// @Produce json
// @Param modelSizeId path int true "Model Size Id"
// @Router /api/wish/cart/reduce/{modelSizeId} [patch]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) reduceNumber(ctx *fiber.Ctx) error {
	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	modelSizeId, err := ctx.ParamsInt("modelSizeId")
	if err != nil {
		appErr := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}
	ex = wh.service.ReduceNumber(ctx.Context(), user.UserId, modelSizeId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Delete several items from cart
// @Security BearerToken
// @Description Delete several items from cart
// @Tags wish
// @Accept json
// @Produce json
// @Param ids path string true "Model Sizes Ids example:([1,2,3])"
// @Router /api/wish/cart/several/{ids} [delete]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (ws *WishHandler) removeSeveralItems(ctx *fiber.Ctx) error {
	claims, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	idsString := ctx.Params("ids")

	if idsString == "" {
		ex := fall.NewErr(fall.VALIDATION_ID, fall.STATUS_BAD_REQUEST)
		return ctx.Status(ex.Status()).JSON(ex)
	}

	idsSlice := strings.Split(idsString, ",")

	var ids []int

	for _, v := range idsSlice {
		id, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	ex = ws.service.RemoveSeveralItems(ctx.Context(), claims.UserId, ids)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}
	resp := fall.GetOk()
	return ctx.Status(resp.Status()).JSON(resp)
}

// @Summary Add item to cart
// @Security BearerToken
// @Description Add item to cart
// @Tags wish
// @Accept json
// @Produce json
// @Param dto body model.AddToCartDto true "Add to cart with dto"
// @Router /api/wish/cart [post]
// @Success 200 {object} fall.AppErr
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (wh *WishHandler) addToCart(ctx *fiber.Ctx) error {

	user, ex := utils.GetLocalSession(ctx)

	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	dto := model.AddToCartDto{}

	err := ctx.BodyParser(&dto)

	if err != nil {
		appErr := fall.NewErr(fall.INVALID_BODY, fall.STATUS_BAD_REQUEST)
		return ctx.Status(appErr.Status()).JSON(appErr)
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := fall.ValidationMessages(error_messages)
		validError := fall.NewValidErr(items)

		return ctx.Status(validError.Status).JSON(validError)
	}

	ex = wh.service.AddToCart(ctx.Context(), dto, user.UserId)
	if ex != nil {
		return ctx.Status(ex.Status()).JSON(ex)
	}

	resp := fall.GetCreated()
	return ctx.Status(resp.Status()).JSON(resp)
}
