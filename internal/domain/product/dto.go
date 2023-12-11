package product

type CreateProductDto struct {
	Title       string  `json:"title" example:"Куртка теплая" validate:"required,min=3"`
	Description *string `json:"description" example:"Отлчиная куртка теплая" validate:"min=10,omitempty"`
	CategoryID  int     `json:"category_id" example:"10" validate:"required,min=1"`
	BrandID     int     `json:"brand_id" example:"10" validate:"required,min=1"`
}

type CreateProductModelDto struct {
	Price     int32  `json:"price" example:"15000" validate:"required,min=1"`
	Discount  *byte  `json:"discount" example:"10"`
	ImagePath string `json:"image_path" validate:"filepath"`
	ProductId int    `json:"product_id" example:"10" validate:"required,min=1"`
}

type CreateProducModelImg struct {
	ImgPath        string `json:"img_path" validate:"required"`
	ProductModelId int    `json:"product_model_id" validate:"required,min=1"`
}

type UpdateProductDto struct {
	Title       *string `json:"title" example:"Куртка теплая" validate:"omitempty"`
	Description *string `json:"description" example:"Отлчиная куртка теплая" validate:"omitempty"`
}

type UpdateProductModelDto struct {
	Price     *int32  `json:"price" example:"15000" validate:"omitempty"`
	Discount  *byte   `json:"discount" example:"10" validate:"omitempty"`
	ImagePath *string `json:"image_path" validate:"omitempty,filepath"`
}
