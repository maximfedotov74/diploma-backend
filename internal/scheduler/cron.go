package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/maximfedotov74/fiber-psql/internal/service"
)

type CronScheduler struct {
	services *service.Services
	cron     *gocron.Scheduler
}

func New(s *service.Services) *CronScheduler {
	cron := gocron.NewScheduler(time.UTC)

	return &CronScheduler{services: s, cron: cron}
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
