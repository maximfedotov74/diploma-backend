package handler

import "github.com/gofiber/fiber/v2"

type DeliveryService interface {
	Create()
}

type DeliveryHandler struct {
	service DeliveryService
}

func NewDeliveryHandler(service DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: service}
}

// @Summary Update category
// @Description Update category
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "id parameter"
// @Router /api/category/{id} [post]
// @Success 200
// @Failure 400 {object} fall.ValidationError
// @Failure 404 {object} fall.AppErr
// @Failure 500 {object} fall.AppErr
func (h *DeliveryHandler) Create(ctx *fiber.Ctx) error {
	h.service.Create()
	return nil
}
