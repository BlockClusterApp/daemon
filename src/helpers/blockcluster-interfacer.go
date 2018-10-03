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
	Valid bool
}

type LicenceValidationResponse struct {
	Success bool `json:"success"`
	Token string `json:"message"`
	Error string `json:"error"`
	ErrorCode int `json:"errorCode"`
}


var Blockcluster BlockClusterType

var BASE_URL = "https://enterprise.blockcluster.io"
//var BASE_URL = "https://a1049eab.ngrok.io"


func (bc *BlockClusterType) SendRequest(path string, body string) (string,error) {
	url := fmt.Sprintf("%s%s", BASE_URL, path)
	res, err := MakeHTTPRequest(url, http.MethodPost, body)
	return res,err
}


// Duplicate function to account for cyclic import
func (bc *BlockClusterType) Reauthorize() {
	path := "/licence/validate"
	jsonBody:= fmt.Sprintf(`{"licence": "%s"}`, base64.StdEncoding.EncodeToString([]byte(bc.Licence.Key)))

	res,err := bc.SendRequest(path, jsonBody)

	if err != nil {
		return
	}

	var licenceResponse = &LicenceValidationResponse{}
	err = json.Unmarshal([]byte(res), licenceResponse)

	if err != nil {
		GetLogger().Printf("Error parsing response %s", err.Error())
		return
	}

	bc.AuthToken = licenceResponse.Token
	bc.Licence.Key = GetLicence().Key
}

func GetBlockclusterInstance() *BlockClusterType {
	return &Blockcluster
}
