package model

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

type CreateProductDto struct {
	Title       string  `json:"title" example:"Куртка теплая" validate:"required,min=3"`
	Description *string `json:"description" example:"Отлчиная куртка теплая" validate:"min=10,omitempty"`
	CategoryId  int     `json:"category_id" example:"10" validate:"required,min=1"`
	BrandId     int     `json:"brand_id" example:"10" validate:"required,min=1"`
}

type CreateProductModelDto struct {
	Price     int32  `json:"price" example:"15000" validate:"required,min=1"`
	Discount  *byte  `json:"discount" example:"10"`
	ImagePath string `json:"image_path" validate:"filepath"`
	ProductId int    `json:"product_id" example:"10" validate:"required,min=1"`
}

type Product struct {
	Id          int           `json:"id" validate:"required" example:"1"`
	Title       string        `json:"title" validate:"required" example:"Куртка"`
	Description *string       `json:"description" example:"описание товара..."`
	Category    CategoryModel `json:"category" validate:"required"`
	Brand       Brand         `json:"brand" validate:"required"`
}

type ProductRelation struct {
	Id           int                  `json:"id" validate:"required" example:"1"`
	Title        string               `json:"title" validate:"required" example:"Куртка"`
	Description  *string              `json:"description" example:"описание товара..."`
	Category     CategoryModel        `json:"category" validate:"required"`
	Brand        Brand                `json:"brand" validate:"required"`
	CurrentModel ProductModelRelation `json:"model"`
}

type ProductModel struct {
	Id        int    `json:"id" example:"1" validate:"required"`
	Price     int32  `json:"price" example:"15000" validate:"required"`
	Discount  *byte  `json:"discount"`
	ProductId int    `json:"product_id" validate:"required"`
	Slug      string `json:"slug" validate:"required"`
	Article   string `json:"article" validate:"required,min=1"`
	ImagePath string `json:"image_path" validate:"required"`
}

type ProductModelRelation struct {
	Id        *int                  `json:"id" example:"1"`
	Price     *int32                `json:"price" example:"15000"`
	Discount  *byte                 `json:"discount"`
	ProductId *int                  `json:"product_id"`
	Slug      string                `json:"slug" validate:"required"`
	Article   string                `json:"article" validate:"required,min=1"`
	ImagePath string                `json:"image_path" validate:"required"`
	Sizes     []ProductModelSize    `json:"sizes" validate:"required"`
	Options   []*ProductModelOption `json:"options"`
	Images    []ProductModelImg     `json:"images" validate:"required"`
}

type ProductModelSize struct {
	SizeId      int    `json:"size_id" example:"1" validate:"required"`
	ModelId     int    `json:"model_id" example:"2" validate:"required"`
	SizeModelId int    `json:"size_model_id" example:"3" validate:"required"`
	Literal     string `json:"literal" example:"M" validate:"required"`
	Value       string `json:"size_value" example:"44" validate:"required"`
	InStock     int    `json:"in_stock" example:"120" validate:"required"`
}

type ProductModelOption struct {
	Id             int                       `json:"id" example:"4" validate:"required"`
	Title          string                    `json:"title" example:"Цвет" validate:"required"`
	Slug           string                    `json:"slug" example:"color" validate:"required"`
	ProductModelId int                       `json:"-" example:"4" validate:"required"`
	Values         []ProductModelOptionValue `json:"values"`
}

type ProductModelOptionValue struct {
	Id             int     `json:"id" example:"44" validate:"required"`
	Value          string  `json:"value" example:"Желтый" validate:"required"`
	Info           *string `json:"info"`
	OptionId       int     `json:"-" example:"4" validate:"required"`
	ProductModelId int     `json:"-" example:"4" validate:"required"`
}

type ProductModelImg struct {
	Id             int    `json:"id" example:"1" validate:"required"`
	ImgPath        string `json:"img_path" validate:"required"`
	ProductModelId int    `json:"-" validate:"required"`
}

type Color struct {
	ModelId *int    `json:"-" validate:"required"`
	Value   *string `json:"value" validate:"required"`
}

type ProductModelColors struct {
	Id    *int    `json:"id" example:"1" validate:"required"`
	Slug  *string `json:"slug" validate:"required"`
	Image string  `json:"image_path" validate:"required"`
	Color Color   `json:"color" validate:"required"`
}

type AdminProductResponse struct {
	Products []*AdminProduct `json:"products"`
	Total    int             `json:"total" validate:"required"`
}

type AdminProduct struct {
	Id          int                         `json:"id" example:"1" validate:"required"`
	Title       string                      `json:"title" example:"Куртка теплая" validate:"required"`
	Description *string                     `json:"description" example:"Отлчиная куртка теплая"`
	Category    CategoryModel               `json:"category" validate:"required"`
	Brand       Brand                       `json:"brand" validate:"required"`
	Models      []AdminProductModelRelation `json:"models" validate:"required"`
}

type AdminProductModelRelation struct {
	Id        *int    `json:"id" example:"1"`
	Price     *int32  `json:"price" example:"15000"`
	Discount  *byte   `json:"discount"`
	Slug      string  `json:"slug" validate:"required"`
	Article   string  `json:"article" validate:"required,min=1"`
	ImagePath *string `json:"image_path"`
	ProductId *int    `json:"product_id"`
}
