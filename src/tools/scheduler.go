package tools

import (
	"github.com/BlockClusterApp/daemon/src/tools/tasks"
	"github.com/jasonlvhit/gocron"
	"log"
)

func StartScheduler() {
	log.Println("Starting gocron")
	gocron.Start()

	gocron.Every(10).Seconds().Do(tasks.FetchNodeInformation)
}