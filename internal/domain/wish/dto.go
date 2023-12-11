package wish

type AddToCartDto struct {
	ModelSizeId int `json:"model_size_id" validate:"required,min=1"`
}

type AddToWishDto struct {
	ModelId int `json:"model_id" validate:"required,min=1"`
}
