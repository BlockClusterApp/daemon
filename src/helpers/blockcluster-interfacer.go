package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type BlockClusterType struct {
	Licence LicenceConfig
	AuthToken string
}

type LicenceValidationResponse struct {
	Success bool `json:"success"`
	Token string `json:"message"`
	Error string `json:"error"`
}

var Blockcluster BlockClusterType

var BASE_URL = "https://enterprise.blockcluster.io"
//var BASE_URL = "https://b9673448.ngrok.io"

func (bc *BlockClusterType) FetchLicenceDetails() {
	url := fmt.Sprintf("%s/licence/validate", BASE_URL)

	jsonBody:= fmt.Sprintf(`{"licence": "%s"}`, base64.StdEncoding.EncodeToString([]byte(Blockcluster.Licence.Key)))

	//log.Println("Request Payload", jsonBody);
	res, err := MakeHTTPRequest(url, http.MethodPost, jsonBody)

	if err != nil {
		return
	}

	var licenceResponse = &LicenceValidationResponse{}
	err = json.Unmarshal([]byte(res), licenceResponse)

	if err != nil {
		GetLogger().Printf("Error parsing response %s", err.Error())
		return
	}

	//log.Printf("Licence response %s", res)
	bc.AuthToken = licenceResponse.Token
	Blockcluster.AuthToken = licenceResponse.Token
}

func GetBlockclusterInstance() *BlockClusterType {
	return &Blockcluster
}

func UpdateBlockClusterInstance(bc BlockClusterType)  {
	Blockcluster = bc
}