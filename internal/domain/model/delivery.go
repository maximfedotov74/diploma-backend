package model

type DeliveryPoint struct {
	Id           int     `json:"delivery_point_id" validate:"required"`
	Title        string  `json:"title" validate:"required"`
	City         string  `json:"city" validate:"required"`
	Address      string  `json:"address" validate:"required"`
	WithFitting  bool    `json:"with_fitting" validate:"required"`
	WorkSchedule string  `json:"work_schedule" validate:"required"`
	Coords       string  `json:"coords" validate:"required"`
	Info         *string `json:"info"`
}

type CreateDeliveryPointDto struct {
	Title        string  `json:"title" validate:"required,min=2"`
	City         string  `json:"city" validate:"required,min=2"`
	Address      string  `json:"address" validate:"required,min=15"`
	WithFitting  bool    `json:"with_fitting" validate:"boolean"`
	WorkSchedule string  `json:"work_schedule" validate:"required,min=2"`
	Coords       string  `json:"coords" validate:"required,min=4"`
	Info         *string `json:"info" validate:"omitempty,min=4"`
}

type UpdateDeliveryPointDto struct {
	Title        *string `json:"title" validate:"omitempty,min=2"`
	City         *string `json:"city" validate:"omitempty,min=2"`
	Address      *string `json:"address" validate:"omitempty,min=15"`
	WithFitting  *bool   `json:"with_fitting" validate:"omitempty,boolean"`
	WorkSchedule *string `json:"work_schedule" validate:"omitempty,min=2"`
	Coords       *string `json:"coords" validate:"omitempty,min=4"`
	Info         *string `json:"info" validate:"omitempty,min=4"`
}
