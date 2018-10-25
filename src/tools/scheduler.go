package tools

import (
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/tools/tasks"
	"github.com/jasonlvhit/gocron"
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

	tasks.UpdateHyperionPorts()

	gocron.Every(10).Minutes().Do(tasks.FetchNodeInformation)
	gocron.Every(5).Minutes().Do(tasks.ValidateLicence)
	gocron.Every(5).Minutes().Do(tasks.FetchPodInformation)
	gocron.Every(5).Hours().Do(tasks.RefreshImagePullSecrets)
	gocron.Every(2).Minutes().Do(tasks.UpdateHyperionPorts)
}
