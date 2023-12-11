package wish

import (
	"strconv"
	"strings"

	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type Repository interface {
	AddToCart(modelSizeId int, userId int) exception.Error
	FindModelInUserCart(modelSizeId int, userId int) (*CartItemModel, exception.Error)
	DeleteFromCart(cartItemId int) exception.Error
	UpdateCartItem(cartItemId int, newQuantity int) exception.Error
	RemoveSeveralItems(cartIds string) exception.Error
	GetUserCart(userId int) ([]CartItem, exception.Error)
	DeleteFromWish(wishId int) exception.Error
	FindWishItem(modelId int, userId int) (*int, exception.Error)
	AddToWish(modelId int, userId int) exception.Error
}

type WishService struct {
	repo Repository
}

func NewWishService(repo Repository) *WishService {
	return &WishService{repo: repo}
}

func (ws *WishService) AddToCart(dto AddToCartDto, userId int) exception.Error {
	item, _ := ws.repo.FindModelInUserCart(dto.ModelSizeId, userId)
	if item != nil {
		return exception.NewErr(CART_ITEM_ALREADY_IN_CART, exception.STATUS_BAD_REQUEST)
	}
	err := ws.repo.AddToCart(dto.ModelSizeId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WishService) DeleteFromCart(userId int, modelSizeId int) exception.Error {
	item, ex := ws.repo.FindModelInUserCart(modelSizeId, userId)

	if ex != nil {
		return ex
	}

	ex = ws.repo.DeleteFromCart(item.CartItemId)
	if ex != nil {
		return ex
	}
	return nil
}

func (ws *WishService) IncreaseNumber(userId int, modelSizeId int) exception.Error {
	item, ex := ws.repo.FindModelInUserCart(modelSizeId, userId)

	if ex != nil {
		return ex
	}

	newQuantity := item.Quantity + 1

	if item.InStock < newQuantity {
		exception.NewErr(QUANTITY_MORE_THAN_IN_STOCK, exception.STATUS_BAD_REQUEST)
	}

	ex = ws.repo.UpdateCartItem(item.CartItemId, newQuantity)
	if ex != nil {
		return ex
	}
	return nil
}

func (ws *WishService) ReduceNumber(userId int, modelSizeId int) exception.Error {
	item, ex := ws.repo.FindModelInUserCart(modelSizeId, userId)

	if ex != nil {
		return ex
	}

	newQuantity := item.Quantity - 1

	if newQuantity <= 0 {
		return exception.NewErr(QUANTITY_LESS_THAN_IN_ZERO, exception.STATUS_BAD_REQUEST)
	}

	ex = ws.repo.UpdateCartItem(item.CartItemId, newQuantity)
	if ex != nil {
		return ex
	}
	return nil
}

func (ws *WishService) RemoveSeveralItems(userId int, modelSizesIds []int) exception.Error {
	var cartIds []string

	for id := range modelSizesIds {
		item, ex := ws.repo.FindModelInUserCart(id, userId)
		if ex != nil {
			continue
		}
		cartIds = append(cartIds, strconv.Itoa(item.CartItemId))
	}

	if len(cartIds) > 0 {
		ex := ws.repo.RemoveSeveralItems(strings.Join(cartIds, ","))
		if ex != nil {
			return ex
		}
		return nil
	}

	return exception.NewErr(CART_ITEM_NOT_FOUND, exception.STATUS_NOT_FOUND)

}

func (ws *WishService) GetUserCart(userId int) ([]CartItem, exception.Error) {
	items, err := ws.repo.GetUserCart(userId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (ws *WishService) ToggleWish(modelId int, userId int) exception.Error {
	existId, _ := ws.repo.FindWishItem(modelId, userId)
	if existId != nil {
		ex := ws.repo.DeleteFromWish(*existId)
		if ex != nil {
			return ex
		}
		return nil
	}
	ex := ws.repo.AddToWish(modelId, userId)
	if ex != nil {
		return ex
	}
	return nil
}
