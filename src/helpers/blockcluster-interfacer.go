package helpers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

type BlockCluster struct {
	Licence LicenceConfig
}

var BASE_URL = "https://enterprise.blockcluster.io"

func (bc BlockCluster) FetchLicenceDetails() {
	url := fmt.Sprintf("%s/licence/validate", BASE_URL)

	jsonBody:= fmt.Sprintf(`{"key": "%s"}`, base64.StdEncoding.EncodeToString([]byte(bc.Licence.Key)))

	res, err := MakeHTTPRequest(url, http.MethodPost, jsonBody)

	if err != nil {
		return
	}

	log.Printf("Licence response %s", res);
}