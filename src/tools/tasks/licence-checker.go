package tasks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"time"
)

func updateWebAppDeployment() {

	webAppPods := helpers.FetchPod("app%3Dblockcluster-app")
	if len(webAppPods.Items) == 0 {
		return
	}
	for i := 0; i < len(webAppPods.Items); i++ {
		go func(i int) {
			// Don't delete all the pods at the same time.
			sleepDuration := time.Duration(i * 60)
			time.Sleep(sleepDuration * time.Second)
			var pod = webAppPods.Items[i]
			helpers.DeletePod(pod.Metadata.Namespace, pod.Metadata.Name)
		}(i)
	}

}

func handleVersionMetadata(licenceResponse *dtos.LicenceValidationResponse) {
	bc := helpers.GetBlockclusterInstance()
	if licenceResponse.Metadata.BlockClusterAgentVersion != helpers.CURRENT_AGENT_VERSION {
		// delete this pod so that it can fetch new image
		blockClusterPods := helpers.FetchPod("app%3Dblockcluster-agent")
		if len(blockClusterPods.Items) == 0 {
			return
		}
		for i := 0; i < len(blockClusterPods.Items); i++ {
			go func(i int) {
				// Don't delete all the pods at the same time.
				sleepDuration := time.Duration(i * 20)
				time.Sleep(sleepDuration * time.Second)
				var pod = blockClusterPods.Items[i]
				helpers.DeletePod(pod.Metadata.Namespace, pod.Metadata.Name)
			}(i)
		}
	}


	webAppMeta := helpers.GetCurrentWebAppMeta()

	if licenceResponse.Metadata.WebAppVersion != "" && licenceResponse.Metadata.WebAppVersion != webAppMeta.WebAppVersion {
		if bc.Metadata.ShouldDaemonDeployWebapp {
			updateWebAppDeployment()
		}
	}

}

func ValidateLicence() {
	helpers.UpdateLicence()
	licence := helpers.GetLicence()
	bc := helpers.GetBlockclusterInstance()
	webAppMeta := helpers.GetCurrentWebAppMeta()
	bc.Licence = licence

	if bc.AgentInfo.WebAppVersion == "" {
		bc.AgentInfo.WebAppVersion = "NotFetched"
	}

	path := "/licence/validate"
	jsonBody := fmt.Sprintf(`{"licence": "%s", "daemonVersion": "%s", "webAppVersion": "%s"}`, base64.StdEncoding.EncodeToString([]byte(licence.Key)), helpers.CURRENT_AGENT_VERSION, webAppMeta.WebAppVersion)

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
