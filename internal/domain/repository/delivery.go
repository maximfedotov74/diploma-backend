package repository

import "fmt"

type DeliveryRepository struct {
	A string
}

func NewDeliveryRepository() *DeliveryRepository {
	return &DeliveryRepository{A: "Hello! Test start"}
}

func (r *DeliveryRepository) Create() {
	fmt.Println(r.A)
}
