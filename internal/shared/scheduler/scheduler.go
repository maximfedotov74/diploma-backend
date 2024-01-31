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

func (s *SchedulerService) Start() {
	s.Cron.StartAsync()
}

func (s *SchedulerService) Shutdown() {
	s.Cron.Stop()
}
