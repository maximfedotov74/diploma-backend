package service

type DeliveryRepository interface {
	Create()
}

type DeliveryService struct {
	repo DeliveryRepository
}

func NewDeliveryService(repo DeliveryRepository) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) Create() {
	s.repo.Create()
}
