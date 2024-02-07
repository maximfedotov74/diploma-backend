package model

type CartItemModel struct {
	CartItemId  int   `json:"cart_item_id" validate:"required"`
	UserId      int   `json:"user_id" validate:"required"`
	ModelSizeId int   `json:"model_size_id" validate:"required"`
	ModelId     int   `json:"model_id" validate:"required"`
	Price       int   `json:"price" validate:"required"`
	Discount    *byte `json:"discount"`
	Quantity    int   `json:"quantity" validate:"required"`
	InStock     int   `json:"in_stock" validate:"required"`
}

type CartItem struct {
	Id        int               `json:"cart_item_id" validate:"required"`
	Quantity  int               `json:"quantity" validate:"required"`
	ModelSize CartItemModelSize `json:"cart_item_model_size" validate:"required"`
}

type CartItemModelSize struct {
	Id           int                  `json:"model_size_id" validate:"required"`
	InStock      int                  `json:"in_stock" validate:"required"`
	LiteralSize  string               `json:"literal_size" validate:"required"`
	SizeId       int                  `json:"size_id" validate:"required"`
	SizeValue    string               `json:"size_value" validate:"required"`
	ProductModel CartItemProductModel `json:"model" validate:"required"`
}

type CartItemProductModel struct {
	Id        int              `json:"model_id" validate:"required"`
	Price     int32            `json:"price" example:"15000" validate:"required"`
	Discount  *byte            `json:"discount"`
	ImagePath string           `json:"image_path" validate:"required"`
	ProductId int              `json:"product_id" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Slug      string           `json:"slug" validate:"required"`
	Category  CartItemCategory `json:"category" validate:"required"`
	Brand     CartItemBrand    `json:"brand" validate:"required"`
}

type CartItemCategory struct {
	Id         int    `json:"category_id" validate:"required"`
	Title      string `json:"title" validate:"required"`
	ShortTitle string `json:"short_title" validate:"required"`
	Slug       string `json:"slug" validate:"required"`
}

type CartItemBrand struct {
	Id    int    `json:"brand_id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Slug  string `json:"slug" validate:"required"`
}

type AddToCartDto struct {
	ModelSizeId int `json:"model_size_id" validate:"required,min=1"`
}

type AddToWishDto struct {
	ModelId int `json:"model_id" validate:"required,min=1"`
}
