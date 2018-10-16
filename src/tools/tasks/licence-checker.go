package tasks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"strings"
	"time"
)

var CURRENT_AGENT_VERSION = "1.0";

type LicenceValidationResponse struct {
	Success bool `json:"success"`
	Token string `json:"message"`
	Error string `json:"error"`
	ErrorCode int `json:"errorCode"`
	Metadata struct {
		BlockClusterAgentVersion string `json:"blockclusterAgentVersion"`
		WebAppVersion string `json:"webappVersion"`
		ShouldDaemonDeployWebapp bool `json:"shouldDaemonDeployWebapp"`
	} `json:"metadata"`
}

func updateWebAppDeployment(newImageTag string) {
	deployment := helpers.FetchDeployment("name%3Dblockcluster-app")
	if deployment == nil {
		return
	}
	webAppIndex := -1
	for i := 0 ; i <  len(deployment.Items[0].Spec.Template.Spec.Containers) ; i++ {
		if deployment.Items[0].Spec.Template.Spec.Containers[i].Name == "blockcluster-webapp" {
			webAppIndex = i
		}
	}

	if webAppIndex < 0 {
		helpers.GetLogger().Printf("Blockcluster-webapp container not found while updating deployment")
		return
	}

	image := deployment.Items[0].Spec.Template.Spec.Containers[webAppIndex].Image
	imageRepo := strings.Split(image, ":")[0]

	deployment.Items[0].Spec.Template.Spec.Containers[webAppIndex].Image = fmt.Sprintf("%s:%s", imageRepo, newImageTag)

	helpers.UpdateDeployment(deployment)
}

func handleVersionMetadata(licenceResponse *LicenceValidationResponse) {
	if licenceResponse.Metadata.BlockClusterAgentVersion != CURRENT_AGENT_VERSION {
		// delete this pod so that it can fetch new image
		blockClusterPods := helpers.FetchPod("app%3Dblockcluster-agent")
		for i := 0 ; i < len(blockClusterPods.Items) ; i++ {
			go func() {
				// Don't delete all the pods at the same time.
				sleepDuration := time.Duration(i * 20)
				time.Sleep( sleepDuration * time.Second)
				var pod = blockClusterPods.Items[i]
				helpers.DeletePod(pod.Metadata.Namespace, pod.Metadata.Name)
			}()
		}
	}

	webAppPodInfo := helpers.FetchPod("app%3Dblockcluster-app")
	if webAppPodInfo == nil {
		return
	}

	var appContainer dtos.Container
	for _, container := range webAppPodInfo.Items[0].Spec.Containers {
		if container.Name == "app" {
			appContainer = container
			break
		}
 	}
	imageTag := (strings.Split(appContainer.Image, ":"))[0]
	if licenceResponse.Metadata.WebAppVersion != "" && licenceResponse.Metadata.WebAppVersion != imageTag {
		updateWebAppDeployment(imageTag)
	}

}

func ValidateLicence() {
	helpers.UpdateLicence()
	licence := helpers.GetLicence()
	bc := helpers.GetBlockclusterInstance()
	bc.Licence = licence

	path := "/licence/validate"
	jsonBody:= fmt.Sprintf(`{"licence": "%s"}`, base64.StdEncoding.EncodeToString([]byte(licence.Key)))

	res,err := bc.SendRequest(path, jsonBody)

	if err != nil {
		return
	}

	var licenceResponse = &LicenceValidationResponse{}
	err = json.Unmarshal([]byte(res), licenceResponse)

	if err != nil {
		helpers.GetLogger().Printf("Error parsing response %s", err.Error())
		return
	}

	bc.AuthToken = licenceResponse.Token
	bc.Licence.Key = helpers.GetLicence().Key

	if licenceResponse.Success != true && licenceResponse.ErrorCode == 401 {
		bc.AuthRetryCount += 1
		if bc.AuthRetryCount >= 3 {
			bc.Valid = false
		}
	} else {
		bc.Valid = true
		bc.AuthRetryCount = 0
	}

	handleVersionMetadata(licenceResponse)
}
