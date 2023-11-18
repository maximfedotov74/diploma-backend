package product

type CreateProductDto struct {
	Title       string  `json:"title" example:"Куртка теплая" validate:"required,min=3"`
	Description *string `json:"description" example:"Отлчиная куртка теплая"`
	CategoryID  int     `json:"category_id" example:"10" validate:"required,min=1"`
	BrandID     int     `json:"brand_id" example:"10" validate:"required,min=1"`
}

type CreateProductModelDto struct {
	Price     float32 `json:"price" example:"15000" validate:"required,min=1"`
	Discount  *byte   `json:"discount" example:"10"`
	ProductId int     `json:"product_id" example:"10" validate:"required,min=1"`
}

type CreateProducModelImg struct {
	ImgPath        string `json:"img_path" validate:"required"`
	Main           *bool  `json:"main"`
	ProductModelId int    `json:"product_model_id" validate:"required,min=1"`
}
