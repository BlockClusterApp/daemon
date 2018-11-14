package tools

import (
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/tools/tasks"
	"github.com/jasonlvhit/gocron"
	"os"
	"time"
)

func StartScheduler() {
	log := helpers.GetLogger()
	log.Println("Starting gocron")
	gocron.Start()

	tasks.ValidateLicence()
	go func() {
		log.Println("Initial image pull secrets")
		time.Sleep(20 * 1000)
		tasks.RefreshImagePullSecrets()
	}()

	tasks.ClearLogFile()

	tasks.UpdateHyperionPorts()
	tasks.CheckHyperionScaler()

	gocron.Every(5).Minutes().Do(tasks.ValidateLicence)

	// The below tasks also updates the cluster config being sent to webapp periodically
	gocron.Every(2).Minutes().Do(tasks.UpdateHyperionPorts)
	gocron.Every(30).Seconds().Do(tasks.UpdateMetrics)

	if os.Getenv("GO_ENV") == "development" {
		return
	}
	gocron.Every(9).Minutes().Do(tasks.FetchNodeInformation)
	gocron.Every(3).Minutes().Do(tasks.FetchPodInformation)
	gocron.Every(5).Hours().Do(tasks.RefreshImagePullSecrets)
	gocron.Every(1).Day().Do(tasks.ClearLogFile)
	gocron.Every(1).Hour().Do(tasks.CheckHyperionScaler)
}
