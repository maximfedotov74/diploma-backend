package brand

type Brand struct {
	Id          int     `json:"id" example:"2"`
	Title       string  `json:"title" example:"adidas"`
	Slug        string  `json:"slug" example:"adidas"`
	Description *string `json:"description" example:"бренд одежды..."`
	ImgPath     *string `json:"img_path" example:"/static/category/adidas.webp"`
}
