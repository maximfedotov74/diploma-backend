package model

type SimilarProductsFilter string

const (
	SimilarByCategory SimilarProductsFilter = "WHERE ct.category_id = $1 AND b.brand_id != $2"
	SimilarByBrand    SimilarProductsFilter = "WHERE b.brand_id = $1 AND ct.category_id = $2"
)

type Views struct {
	Count *int `json:"count"`
}

type CreateProducModelImg struct {
	ImgPath        string `json:"img_path" validate:"required"`
	ProductModelId int    `json:"product_model_id" validate:"required,min=1"`
}

type UpdateProductDto struct {
	Title       *string `json:"title" example:"Куртка теплая" validate:"omitempty,min=3"`
	Description *string `json:"description" example:"Отлчиная куртка теплая" validate:"omitempty,min=10"`
}

type UpdateProductModelDto struct {
	Price     *int32  `json:"price" example:"15000" validate:"omitempty"`
	Discount  *byte   `json:"discount" example:"10" validate:"omitempty"`
	ImagePath *string `json:"image_path" validate:"omitempty,filepath"`
}

type CreateProductDto struct {
	Title       string  `json:"title" example:"Куртка теплая" validate:"required,min=3"`
	Description *string `json:"description" example:"Отлчиная куртка теплая" validate:"omitempty,min=10"`
	CategoryId  int     `json:"category_id" example:"10" validate:"required,min=1"`
	BrandId     int     `json:"brand_id" example:"10" validate:"required,min=1"`
}

type CreateProductModelDto struct {
	Price     int32  `json:"price" example:"15000" validate:"required,min=1"`
	Discount  *byte  `json:"discount" example:"10"`
	ImagePath string `json:"image_path" validate:"required,filepath"`
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
	CurrentModel ProductModelRelation `json:"model" validate:"required"`
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
	Slug      *string               `json:"slug" validate:"required"`
	Article   *string               `json:"article" validate:"required,min=1"`
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

type OrderProductModelSize struct {
	SizeModelId int   `json:"size_model_id" validate:"required"`
	ModelId     int   `json:"model_id" validate:"required"`
	Price       int   `json:"price" validate:"required"`
	Discount    *byte `json:"discount"`
	InStock     int   `json:"in_stock" validate:"required"`
}

type ProductModelOption struct {
	Id                   int                       `json:"id" example:"4" validate:"required"`
	Title                string                    `json:"title" example:"Цвет" validate:"required"`
	Slug                 string                    `json:"slug" example:"color" validate:"required"`
	ProductModelId       int                       `json:"-" example:"4" validate:"required"`
	ProductModelOptionId int                       `json:"pmop_id" example:"4" validate:"required"`
	Values               []ProductModelOptionValue `json:"values"`
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
	Id          int           `json:"id" example:"1" validate:"required"`
	Title       string        `json:"title" example:"Куртка теплая" validate:"required"`
	Description *string       `json:"description" example:"Отлчиная куртка теплая"`
	Category    CategoryModel `json:"category" validate:"required"`
	Brand       Brand         `json:"brand" validate:"required"`
}

type AdminProductModelRelation struct {
	Id        int    `json:"id" example:"1" validate:"required"`
	Price     int32  `json:"price" example:"15000" validate:"required"`
	Discount  *byte  `json:"discount"`
	Slug      string `json:"slug" validate:"required"`
	Article   string `json:"article" validate:"required,min=1"`
	ImagePath string `json:"image_path" validate:"required"`
	ProductId int    `json:"product_id" validate:"required"`
}

type CatalogModelBrand struct {
	Id    int    `json:"id" example:"2" validate:"required"`
	Title string `json:"title" example:"adidas" validate:"required"`
	Slug  string `json:"slug" example:"adidas" validate:"required"`
}

type CatalogModelCategory struct {
	Id         int    `json:"category_id" db:"category_id" example:"3" validate:"required"`
	Title      string `json:"title" db:"title" example:"Верхняя одежда" validate:"required"`
	ShortTitle string `json:"short_title" db:"title" example:"одежда" validate:"required"`
	Slug       string `json:"slug" db:"slug" example:"verhnia-odezhda" validate:"required"`
}

type CatalogProductModel struct {
	ProductId     int                  `json:"product_id" example:"1" validate:"required"`
	Title         string               `json:"product_title" example:"Ботинки" validate:"required"`
	Slug          string               `json:"product_slug" example:"botinki" validate:"required"`
	Article       string               `json:"article" validate:"required,min=1"`
	ModelId       int                  `json:"model_id" example:"1" validate:"required"`
	Price         int                  `json:"model_price" example:"10000" validate:"required"`
	Discount      *byte                `json:"model_discount" example:"15"`
	MainImagePath string               `json:"model_main_image_path" example:"/static/category/test.webp" validate:"required"`
	Brand         CatalogModelBrand    `json:"brand" validate:"required"`
	Category      CatalogModelCategory `json:"category" validate:"required"`
	Images        []*ProductModelImg   `json:"images" validate:"required"`
	Sizes         []*ProductModelSize  `json:"sizes" validate:"required"`
}

type CatalogResponse struct {
	Models     []*CatalogProductModel `json:"models"`
	TotalCount int                    `json:"total_count" example:"100" validate:"required"`
}

type SearchProductModel struct {
	ProductId     int           `json:"product_id" example:"1" validate:"required"`
	Title         string        `json:"product_title" example:"Ботинки" validate:"required"`
	Slug          string        `json:"product_slug" example:"botinki" validate:"required"`
	Article       string        `json:"article" validate:"required,min=1"`
	ModelId       int           `json:"model_id" example:"1" validate:"required"`
	Price         int           `json:"model_price" example:"10000" validate:"required"`
	Discount      *byte         `json:"model_discount" example:"15"`
	MainImagePath string        `json:"model_main_image_path" example:"/static/category/test.webp" validate:"required"`
	Brand         Brand         `json:"brand" validate:"required"`
	Category      CategoryModel `json:"category" validate:"required"`
}
