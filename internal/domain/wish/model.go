package wish

type CartItemModel struct {
	CartItemId  int `json:"cart_item_id"`
	UserId      int `json:"user_id"`
	ModelSizeId int `json:"model_size_id"`
	Quantity    int `json:"quantity"`
	InStock     int `json:"in_stock"`
}

type CartItem struct {
	Id        int               `json:"cart_item_id"`
	Quantity  int               `json:"quantity"`
	ModelSize CartItemModelSize `json:"cart_item_model_size"`
}

type CartItemModelSize struct {
	Id           int                  `json:"model_size_id"`
	InStock      int                  `json:"in_stock"`
	LiteralSize  string               `json:"literal_size"`
	SizeId       int                  `json:"size_id"`
	SizeValue    string               `json:"size_value"`
	ProductModel CartItemProductModel `json:"model"`
}

type CartItemProductModel struct {
	Id        int             `json:"model_id"`
	Price     int32           `json:"price" example:"15000"`
	Discount  *byte           `json:"discount"`
	ImagePath string          `json:"image_path"`
	Product   CartItemProduct `json:"product"`
}

type CartItemProduct struct {
	Id       int              `json:"product_id"`
	Title    string           `json:"title"`
	Slug     string           `json:"slug"`
	Category CartItemCategory `json:"category"`
	Brand    CartItemBrand    `json:"brand"`
}

type CartItemCategory struct {
	Id         int    `json:"category_id"`
	Title      string `json:"title"`
	ShortTitle string `json:"short_title"`
	Slug       string `json:"slug"`
}

type CartItemBrand struct {
	Id    int    `json:"brand_id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
