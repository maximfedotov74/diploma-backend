package scheduler

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/maximfedotov74/diploma-backend/internal/domain/repository"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/payment"
)

type OrderScheduler struct {
	cron    *gocron.Scheduler
	db      db.PostgresClient
	repo    repository.OrderRepository
	payment payment.PaymentService
}

func NewOrderScheduler(cron *gocron.Scheduler, db db.PostgresClient, repo repository.OrderRepository,
	payment payment.PaymentService,
) *OrderScheduler {
	return &OrderScheduler{cron: cron, db: db, repo: repo, payment: payment}
}

func (s *OrderScheduler) Start() {

	ctx := context.Background()

	fmt.Println(ctx)

}

//TODO: implement

func (s *OrderScheduler) CheckOrderPayment(ctx context.Context) {
	s.payment.CheckPayment("")
}
