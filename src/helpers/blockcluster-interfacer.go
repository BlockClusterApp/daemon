package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/getsentry/raven-go"
	"net/http"
)

type BlockClusterType struct {
	Licence        LicenceConfig
	AuthToken      string
	Valid          bool
	AuthRetryCount int8
	Metadata       struct {
		BlockClusterAgentVersion string
		WebAppVersion            string
		ShouldDaemonDeployWebapp bool
		ClientID string
	}
}

var Blockcluster BlockClusterType

var BASE_URL = "https://enterprise-api.blockcluster.io/daemon"
//var BASE_URL = "https://3d7089e8.ngrok.io/daemon"

func (bc *BlockClusterType) SendRequest(path string, body string) (string, error) {
	url := fmt.Sprintf("%s%s", BASE_URL, path)
	res, err := MakeHTTPRequest(url, http.MethodPost, body)
	return res, err
}

// Duplicate function to account for cyclic import
func (bc *BlockClusterType) Reauthorize() {
	path := "/licence/validate"
	jsonBody := fmt.Sprintf(`{"licence": "%s"}`, base64.StdEncoding.EncodeToString([]byte(bc.Licence.Key)))

	res, err := bc.SendRequest(path, jsonBody)

	if err != nil {
		return
	}

	var licenceResponse = &dtos.LicenceValidationResponse{}
	err = json.Unmarshal([]byte(res), licenceResponse)

	if err != nil {
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
		})
		GetLogger().Printf("Error parsing response %s", err.Error())
		return
	}

	if licenceResponse.Error != "" {
		GetLogger().Printf("Error from licence validation %s", licenceResponse.Error)
		bc.AuthToken = ""
		return
	}

	bc.AuthToken = licenceResponse.Token
	bc.Licence.Key = GetLicence().Key

}

func GetBlockclusterInstance() *BlockClusterType {
	return &Blockcluster
}

func (bc *BlockClusterType) GetAWSCreds() *dtos.AWSCredsResponse {
	path := "/aws-creds"

	jsonBody := "{}"

	var awsCredsResponse = &dtos.AWSCredsResponse{}

	res, err := bc.SendRequest(path, jsonBody)

	if err != nil {
		GetLogger().Printf("Error fetching aws creds %s", err.Error())
		return awsCredsResponse
	}

	err = json.Unmarshal([]byte(res), awsCredsResponse)
	if err != nil {
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
		})
		GetLogger().Printf("Error unmarshalling aws creds response %s", err.Error())
		return awsCredsResponse
	}

	return awsCredsResponse
}
