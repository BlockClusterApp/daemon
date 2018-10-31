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

func updateWebAppDeployment(newImageTag string) {
	deployment := helpers.FetchDeployment("app%3Dblockcluster-app")
	if deployment == nil {
		return
	}
	webAppIndex := -1
	for i := 0; i < len(deployment.Items[0].Spec.Template.Spec.Containers); i++ {
		if deployment.Items[0].Spec.Template.Spec.Containers[i].Name == "blockcluster-webapp" || deployment.Items[0].Spec.Template.Spec.Containers[i].Name == "app-deploy" {
			webAppIndex = i
			break
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

func handleVersionMetadata(licenceResponse *dtos.LicenceValidationResponse) {
	bc := helpers.GetBlockclusterInstance()
	if licenceResponse.Metadata.BlockClusterAgentVersion != helpers.CURRENT_AGENT_VERSION {
		// delete this pod so that it can fetch new image
		blockClusterPods := helpers.FetchPod("app%3Dblockcluster-agent")
		for i := 0; i < len(blockClusterPods.Items); i++ {
			go func() {
				// Don't delete all the pods at the same time.
				sleepDuration := time.Duration(i * 20)
				time.Sleep(sleepDuration * time.Second)
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
	imageTag := (strings.Split(appContainer.Image, ":"))[1]

	bc.AgentInfo.WebAppVersion = imageTag

	if licenceResponse.Metadata.WebAppVersion != "" && licenceResponse.Metadata.WebAppVersion != imageTag {
		if bc.Metadata.ShouldDaemonDeployWebapp {
			updateWebAppDeployment(imageTag)
		}
	}

}

func ValidateLicence() {
	helpers.UpdateLicence()
	licence := helpers.GetLicence()
	bc := helpers.GetBlockclusterInstance()
	bc.Licence = licence

	if bc.AgentInfo.WebAppVersion == "" {
		bc.AgentInfo.WebAppVersion = "NotFetched"
	}

	path := "/licence/validate"
	jsonBody := fmt.Sprintf(`{"licence": "%s", "daemonVersion": "%s", "webAppVersion": "%s"}`, base64.StdEncoding.EncodeToString([]byte(licence.Key)), helpers.CURRENT_AGENT_VERSION, bc.AgentInfo.WebAppVersion)

	res, err := bc.SendRequest(path, jsonBody)

	if err != nil {
		return
	}

	var licenceResponse = &dtos.LicenceValidationResponse{}
	err = json.Unmarshal([]byte(res), licenceResponse)

	if err != nil {
		helpers.GetLogger().Printf("Error parsing response %s", err.Error())
		return
	}

	helpers.GetLogger().Printf("Got response %s", res)

	bc.AuthToken = licenceResponse.Token
	bc.Licence.Key = helpers.GetLicence().Key

	bc.Metadata = licenceResponse.Metadata

	helpers.RefreshLogger()

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
