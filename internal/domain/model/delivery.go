package model

import "fmt"

type DeliveryRepository struct {
	A string
}

func NewDeliveryRepository() *DeliveryRepository {
	return &DeliveryRepository{}
}

func (r *DeliveryRepository) Create() {
	fmt.Println(r.A)
}
