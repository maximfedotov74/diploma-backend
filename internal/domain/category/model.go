package category

type CatalogCategory struct {
	Id            int               `json:"category_id" db:"category_id" example:"3"`
	Title         string            `json:"title" db:"title" example:"Верхняя одежда"`
	Slug          string            `json:"slug" db:"slug" example:"verhnia-odezhda"`
	ParentId      int               `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Count         int               `json:"count" example:"4"`
	Subcategories []CatalogCategory `json:"subcategories"`
}

type RecursiveCategory struct {
	Id            int                 `json:"category_id" db:"category_id" example:"3"`
	Title         string              `json:"title" db:"title" example:"Верхняя одежда"`
	Slug          string              `json:"slug" db:"slug" example:"verhnia-odezhda"`
	ShortTitle    string              `json:"short_title"`
	ImgPath       *string             `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int                `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Level         uint8               `json:"level"`
	Subcategories []RecursiveCategory `json:"subcategories"`
}
