package scheduler

import (
	"github.com/go-co-op/gocron"
)

type SchedulerService struct {
	cron *gocron.Scheduler
}

func New(cron *gocron.Scheduler) *SchedulerService {

	return &SchedulerService{cron: cron}
}

func (cs *SchedulerService) Start() {
	// cs.cron.Every(5).Seconds().Do(func() {
	// 	fmt.Println("asasdasd")
	// })

	cs.cron.StartAsync()
}

func (cs *SchedulerService) Shutdown() {
	cs.cron.Stop()
}
