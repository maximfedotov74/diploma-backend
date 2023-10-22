package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/maximfedotov74/fiber-psql/internal/handler"
)

type CronScheduler struct {
	handler *handler.Handler
	cron    *gocron.Scheduler
}

func New(h *handler.Handler) *CronScheduler {
	cron := gocron.NewScheduler(time.UTC)

	return &CronScheduler{handler: h, cron: cron}
}

func (cs *CronScheduler) Start() {
	// cs.cron.Every(5).Seconds().Do(func() {
	// 	fmt.Println("asasdasd")
	// })

	cs.cron.StartAsync()
}

func (cs *CronScheduler) Shutdown() {
	cs.cron.Stop()
}
