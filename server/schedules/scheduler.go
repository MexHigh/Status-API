package schedules

import (
	"time"

	"github.com/go-co-op/gocron"
)

// this is required to be package global as
// there can only be one scheduler in a routine
var scheduler *gocron.Scheduler

func init() {
	scheduler = gocron.NewScheduler(time.Now().Location())
	scheduler.StartAsync() // jobs can be registered of the scheduler is already running
}
