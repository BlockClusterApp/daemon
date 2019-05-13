package tools

import (
	"github.com/BlockClusterApp/daemon/src/tools/tasks"
	"github.com/jasonlvhit/gocron"
	"log"
	"os"
	"time"
)

func StartScheduler() {
	log.Println("Starting jobs")
	gocron.Start()
	sleepDuration := time.Duration(1)
	time.Sleep(sleepDuration * time.Second)

	tasks.ValidateLicence()

	tasks.UpdateConfigs()

	if os.Getenv("GO_ENV") == "development" {
		return
	}

	go func() {
		log.Println("Pulling Image Secrets")
		time.Sleep(time.Duration(20) * time.Second)
		tasks.RefreshImagePullSecrets()
	}()

	tasks.ClearLogFile()
	tasks.UpdateConfigs()
	tasks.CheckAllSystems()

	tasks.CheckHyperionScaler()

	gocron.Every(5).Minutes().Do(tasks.ValidateLicence)
	gocron.Every(10).Minutes().Do(tasks.CheckAllSystems)

	gocron.Every(2).Minutes().Do(tasks.UpdateConfigs)
	gocron.Every(30).Seconds().Do(tasks.UpdateMetrics)
	gocron.Every(9).Minutes().Do(tasks.FetchNodeInformation)
	gocron.Every(3).Minutes().Do(tasks.FetchPodInformation)
	gocron.Every(4).Hours().Do(tasks.RefreshImagePullSecrets)
	gocron.Every(1).Day().Do(tasks.ClearLogFile)
	gocron.Every(1).Hour().Do(tasks.CheckHyperionScaler)
}
