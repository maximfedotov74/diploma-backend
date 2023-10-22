package handler

import (
	"github.com/go-co-op/gocron"
)

func (h *Handler) SetupCronJobs(cron *gocron.Scheduler) {

	// cron.Every(5).Seconds().Do(func() {
	// 	fmt.Println("assss")
	// })

	cron.StartAsync()

}
