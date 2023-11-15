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
	Id          int             `json:"id" example:"1"`
	Title       string          `json:"title" example:"Куртка теплая"`
	Slug        string          `json:"slug"`
	Description *string         `json:"description" example:"Отлчиная куртка теплая"`
	Category    ProductCategory `json:"category"`
	Brand       ProductBrand    `json:"brand"`
	Images      []ProductImg    `json:"images"`
}

type ProductModel struct {
	Id        int     `json:"id" example:"1"`
	Price     float32 `json:"price" example:"15000"`
	Discount  *byte   `json:"discount"`
	ProductId int     `json:"product_id"`
}

type ProductImg struct {
	Id      *int    `json:"id" example:"1"`
	ImgPath *string `json:"img_path"`
	Main    *bool   `json:"main"`
}
