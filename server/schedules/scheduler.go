package schedules

import (
	"time"

	"github.com/go-co-op/gocron"
)

var scheduler *gocron.Scheduler

func init() {
	scheduler = gocron.NewScheduler(time.Now().Location())
	scheduler.StartAsync()
}
