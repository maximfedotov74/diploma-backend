package scheduler

import (
	"context"
	"log"

	"github.com/go-co-op/gocron"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/payment"
)

type OrderScheduler struct {
	cron    *gocron.Scheduler
	db      db.PostgresClient
	payment *payment.PaymentService
}

func NewOrderScheduler(cron *gocron.Scheduler, db db.PostgresClient,
	payment *payment.PaymentService,
) *OrderScheduler {
	return &OrderScheduler{cron: cron, db: db, payment: payment}
}

func (s *OrderScheduler) Start() {

	ctx := context.Background()

	go s.CheckOrderPayment(ctx)

}

//TODO: implement

type OP struct {
	OrderId   string
	PaymentId string
}

func (s *OrderScheduler) CheckOrderPayment(ctx context.Context) {

	s.cron.Every(3).Minute().Do(func() {
		q := "SELECT order_id, payment_id FROM public.order WHERE order_payment_method = $1 AND order_status = $2 AND payment_id IS NOT NULL;"

		rows, err := s.db.Query(ctx, q, model.Online, model.WaitingForPayment)

		if err != nil {
			log.Println(err.Error())
			return
		}
		defer rows.Close()

		var ops []OP

		for rows.Next() {
			var orderId string
			var payemntId string
			err := rows.Scan(&orderId, &payemntId)

			if err != nil {
				log.Println(err.Error())
				return
			}
			ops = append(ops, OP{OrderId: orderId, PaymentId: payemntId})
		}

		for _, item := range ops {
			p, err := s.payment.CheckPayment(item.PaymentId)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if p.Status == "succeeded" {
				q := "UPDATE public.order SET order_status = $1 WHERE order_id = $2;"
				_, err := s.db.Exec(ctx, q, model.Paid, item.OrderId)
				if err != nil {
					log.Println(err.Error())
					continue
				}
			}
		}
		log.Println("Order payment checker successfully completed!")
	})
}
