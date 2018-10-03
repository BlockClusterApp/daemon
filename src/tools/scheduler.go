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

	tasks.ValidateLicence()
	gocron.Every(10).Minutes().Do(tasks.FetchNodeInformation)
	gocron.Every(5).Minutes().Do(tasks.ValidateLicence)
	gocron.Every(5).Minutes().Do(tasks.FetchPodInformation)
}