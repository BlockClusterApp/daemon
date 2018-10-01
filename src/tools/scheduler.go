package tools

import (
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/tools/tasks"
	"github.com/jasonlvhit/gocron"
)

func StartScheduler() {
	log := helpers.GetLogger()
	log.Println("Starting gocron")
	gocron.Start()

	gocron.Every(1).Minute().Do(tasks.FetchNodeInformation)
	//gocron.Every(5).Seconds().Do(tasks.ClearLogFile)
	gocron.Every(5).Seconds().Do(tasks.ValidateLicence)
}