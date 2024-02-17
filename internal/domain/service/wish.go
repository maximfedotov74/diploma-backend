package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type wishRepository interface {
	FindModelInUserCart(ctx context.Context, modelSizeId int, userId int) (*model.CartItemModel, fall.Error)
	AddToCart(ctx context.Context, modelSizeId int, userId int) fall.Error
	DeleteFromCart(ctx context.Context, cartItemId int) fall.Error
	UpdateCartItem(ctx context.Context, cartItemId int, newQuantity int) fall.Error
	RemoveSeveralItems(ctx context.Context, tx db.Transaction, cartIds []int) fall.Error
	GetUserCart(ctx context.Context, userId int) ([]model.CartItem, fall.Error)
	AddToWish(ctx context.Context, modelId int, userId int) fall.Error
	FindWishItem(ctx context.Context, modelId int, userId int) (*int, fall.Error)
	DeleteFromWish(ctx context.Context, wishId int) fall.Error
	GetUserWish(ctx context.Context, userId int) ([]*model.CatalogProductModel, fall.Error)
}

type WishService struct {
	repo wishRepository
}

func NewWishService(repo wishRepository) *WishService {
	return &WishService{repo: repo}
}

func (s *WishService) GetUserWish(ctx context.Context, userId int) ([]*model.CatalogProductModel, fall.Error) {
	return s.repo.GetUserWish(ctx, userId)
}

func (s *WishService) FindModelInUserCart(ctx context.Context, modelSizeId int, userId int) (*model.CartItemModel, fall.Error) {
	return s.repo.FindModelInUserCart(ctx, modelSizeId, userId)
}

func (s *WishService) AddToCart(ctx context.Context, dto model.AddToCartDto, userId int) fall.Error {
	item, _ := s.FindModelInUserCart(ctx, dto.ModelSizeId, userId)
	if item != nil {
		return fall.NewErr(msg.CartItemAlreadyInCart, fall.STATUS_BAD_REQUEST)
	}
	return s.repo.AddToCart(ctx, dto.ModelSizeId, userId)
}

func (s *WishService) DeleteFromCart(ctx context.Context, userId int, modelSizeId int) fall.Error {
	item, ex := s.FindModelInUserCart(ctx, modelSizeId, userId)

	if ex != nil {
		return ex
	}

	return s.repo.DeleteFromCart(ctx, item.CartItemId)

}

func (s *WishService) IncreaseNumber(ctx context.Context, userId int, modelSizeId int) fall.Error {
	item, ex := s.FindModelInUserCart(ctx, modelSizeId, userId)

	if ex != nil {
		return ex
	}

	newQuantity := item.Quantity + 1

	if item.InStock < newQuantity {
		return fall.NewErr(msg.QuantityMoreThanInStock, fall.STATUS_BAD_REQUEST)
	}

	return s.repo.UpdateCartItem(ctx, item.CartItemId, newQuantity)

}

func (s *WishService) ReduceNumber(ctx context.Context, userId int, modelSizeId int) fall.Error {
	item, ex := s.FindModelInUserCart(ctx, modelSizeId, userId)

	if ex != nil {
		return ex
	}

	newQuantity := item.Quantity - 1

	if newQuantity <= 0 {
		return fall.NewErr(msg.QuantityLessThanZero, fall.STATUS_BAD_REQUEST)
	}

	return s.repo.UpdateCartItem(ctx, item.CartItemId, newQuantity)
}

func (s *WishService) RemoveSeveralItems(ctx context.Context, userId int, modelSizesIds []int) fall.Error {
	var cartIds []int

	for _, id := range modelSizesIds {
		item, ex := s.FindModelInUserCart(ctx, id, userId)
		if ex != nil {
			continue
		}
		cartIds = append(cartIds, item.CartItemId)
	}

	if len(cartIds) > 0 {
		ex := s.repo.RemoveSeveralItems(ctx, nil, cartIds)
		if ex != nil {
			return ex
		}
		return nil
	}

	return fall.NewErr(msg.CartItemNotFound, fall.STATUS_NOT_FOUND)

}

func (s *WishService) GetUserCart(ctx context.Context, userId int) ([]model.CartItem, fall.Error) {
	return s.repo.GetUserCart(ctx, userId)

}

func (s *WishService) ToggleWish(ctx context.Context, modelId int, userId int) fall.Error {
	existId, _ := s.repo.FindWishItem(ctx, modelId, userId)
	if existId != nil {
		return s.repo.DeleteFromWish(ctx, *existId)
	}
	return s.repo.AddToWish(ctx, modelId, userId)
}
