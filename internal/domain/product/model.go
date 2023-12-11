package product

type ProductCategory struct {
	Id         int     `json:"category_id" db:"category_id" example:"3"`
	Title      string  `json:"title" db:"title" example:"Верхняя одежда"`
	Slug       string  `json:"slug" db:"slug" example:"verhnia-odezhda"`
	ShortTitle string  `json:"short_title"`
	ImgPath    *string `json:"img_path" db:"img_path" example:"/static/example.webp"`
}

type ProductBrand struct {
	Id      int     `json:"id" example:"2"`
	Title   string  `json:"title" example:"adidas"`
	Slug    string  `json:"slug" db:"slug" example:"verhnia-odezhda"`
	ImgPath *string `json:"img_path"`
}

type Product struct {
	Id           int             `json:"id" example:"1"`
	Title        string          `json:"title" example:"Куртка теплая"`
	Slug         string          `json:"slug"`
	Description  *string         `json:"description" example:"Отлчиная куртка теплая"`
	Category     ProductCategory `json:"category"`
	Brand        ProductBrand    `json:"brand"`
	CurrentModel *ProductModel   `json:"current_model"`
}

type ProductWithoutRelations struct {
	Id          int             `json:"id" example:"1"`
	Title       string          `json:"title" example:"Куртка теплая"`
	Slug        string          `json:"slug"`
	Category    ProductCategory `json:"category"`
	Brand       ProductBrand    `json:"brand"`
	Description *string         `json:"description" example:"Отлчиная куртка теплая"`
}

type CatalogModel struct {
}

type AdminProductResponse struct {
	Products []AdminProduct `json:"products"`
	Total    int            `json:"total"`
}

type AdminProduct struct {
	Id          int                         `json:"id" example:"1"`
	Title       string                      `json:"title" example:"Куртка теплая"`
	Slug        string                      `json:"slug"`
	Description *string                     `json:"description" example:"Отлчиная куртка теплая"`
	Category    ProductCategory             `json:"category"`
	Brand       ProductBrand                `json:"brand"`
	Models      []AdminProductModelRelation `json:"models"`
}

type ProductModelWithoutRelations struct {
	Id        int    `json:"id" example:"1"`
	Price     int32  `json:"price" example:"15000"`
	Discount  *byte  `json:"discount"`
	ImagePath string `json:"image_path"`
	ProductId int    `json:"product_id"`
}

type AdminProductModelRelation struct {
	Id        *int    `json:"id" example:"1"`
	Price     *int32  `json:"price" example:"15000"`
	Discount  *byte   `json:"discount"`
	ImagePath *string `json:"image_path"`
	ProductId *int    `json:"product_id"`
}

type ProductModelSize struct {
	SizeId      int    `json:"size_id"`
	ModelId     int    `json:"model_id"`
	SizeModelId int    `json:"size_model_id"`
	Literal     string `json:"literal"`
	InStock     int    `json:"in_stock"`
}

type ProductModel struct {
	Id        *int                 `json:"id" example:"1"`
	Price     *int32               `json:"price" example:"15000"`
	Discount  *byte                `json:"discount"`
	ProductId *int                 `json:"product_id"`
	ImagePath string               `json:"image_path"`
	Sizes     []ProductModelSize   `json:"sizes"`
	Options   []ProductModelOption `json:"options"`
	Images    []ProductModelImg    `json:"images"`
}

type Color struct {
	ModelId *int    `json:"-"`
	Value   *string `json:"value"`
}

type ProductModelColors struct {
	Id          *int    `json:"id" example:"1"`
	ProductSlug *string `json:"slug"`
	Image       string  `json:"image_path"`
	Color       Color   `json:"color"`
}

type ProductModelOption struct {
	Id             *int                      `json:"id" example:"4"`
	Title          *string                   `json:"title" example:"Цвет"`
	Slug           *string                   `json:"slug" example:"color"`
	ProductModelId *int                      `json:"-" example:"4"`
	Values         []ProductModelOptionValue `json:"values"`
}

type ProductModelOptionValue struct {
	Id             *int    `json:"id" example:"44"`
	Value          *string `json:"value" example:"Желтый"`
	Info           *string `json:"info"`
	OptionId       *int    `json:"-" example:"4"`
	ProductModelId *int    `json:"-" example:"4"`
}

type ProductModelImg struct {
	Id             *int    `json:"id" example:"1"`
	ImgPath        *string `json:"img_path"`
	ProductModelId *int    `json:"-"`
}
