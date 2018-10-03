package tasks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
)

type LicenceValidationResponse struct {
	Success bool `json:"success"`
	Token string `json:"message"`
	Error string `json:"error"`
	ErrorCode int `json:"errorCode"`
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
		bc.Valid = false
	} else {
		bc.Valid = true
	}
}
