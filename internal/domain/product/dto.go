package product

type CreateProductDto struct {
	Title       string  `json:"title" example:"Куртка теплая"`
	Description *string `json:"description" example:"Отлчиная куртка теплая"`
	CategoryID  int     `json:"category_id" example:"10"`
	BrandID     int     `json:"brand_id" example:"10"`
}

type CreateProductModelDto struct {
	Price     float32 `json:"price" example:"15000"`
	Discount  *byte   `json:"discount" example:"10"`
	ProductId int     `json:"product_id" example:"10"`
}

type CreateProductImg struct {
	ImgPath   string `json:"img_path"`
	Main      bool   `json:"main"`
	ProductId int    `json:"product_id"`
}
