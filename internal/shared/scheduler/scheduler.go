package scheduler

import (
	"github.com/go-co-op/gocron"
)

type SchedulerService struct {
	Cron *gocron.Scheduler
}

func New(cron *gocron.Scheduler) *SchedulerService {
	return &SchedulerService{Cron: cron}
}

func (cs *SchedulerService) Start() {
	cs.Cron.StartAsync()
}

func (cs *SchedulerService) Shutdown() {
	cs.Cron.Stop()
}
