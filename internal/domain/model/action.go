package model

import "time"

type CreateActionDto struct {
	EndDate     time.Time `json:"end_date" validate:"required,gt"`
	Title       string    `json:"title" validate:"required,min=5"`
	ImgPath     *string   `json:"img_path" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty,min=10"`
}

type AddModelToActionDto struct {
	ProductModelId int    `json:"product_model_id" validate:"required,min=1"`
	ActionId       string `json:"action_id" validate:"required,uuid"`
}

type Action struct {
	Id          string    `json:"id" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Title       string    `json:"title" validate:"required,min=5"`
	IsActivated bool      `json:"is_activated" validate:"required"`
	ImgPath     *string   `json:"img_path" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty,min=10"`
}
