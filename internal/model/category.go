package model

type CategoryType struct {
	Id    int    `json:"category_id" db:"category_id" example:"1"`
	Title string `json:"title" db:"title" example:"Пальто"`
}

type Category struct {
	Id            int        `json:"category_id" db:"category_id" example:"3"`
	Title         string     `json:"title" db:"title" example:"Верхняя одежда"`
	ImgPath       *string    `json:"img_path" db:"img_path" example:"/static/example.webp"`
	ParentId      *int       `json:"parent_category_id" db:"parent_category_id" example:"2"`
	Subcategories []Category `json:"subcategories"`
}

type CreateCategoryDto struct {
	Title    string  `json:"title" db:"title" example:"Верхняя одежда" validate:"required,min=3"`
	ImgPath  *string `json:"img_path" db:"img_path" example:"/static/example.webp" validate:"omitempty"`
	ParentId *int    `json:"parent_category_id" validate:"omitempty"`
}

type CreateCategoryTypeDto struct {
	Title string `json:"title" db:"title" example:"Верхняя одежда" validate:"required,min=3"`
}
