package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"net/http"
)

func FetchPodInformation() {
	log := helpers.GetLogger()
	log.Println("G:TASK Fetching pod information")
	podsInfo, err := helpers.MakeKubeRequest(http.MethodGet, "/api/v1/pods", nil)

	if err != nil {
		return
	}
	NodeMap := &dtos.InfoResponse{}
	err = json.Unmarshal([]byte(podsInfo), NodeMap)

	log.Println(podsInfo)

	if err != nil {
		log.Printf("Error parsing json while fetching node info %s", err.Error())
		return
	}

	requestBody := fmt.Sprintf(`{"info": "%s", "timestamp": "%d"}`, podsInfo, helpers.GetTimeInMillis())
	path := "/info/pods"

	bc := helpers.GetBlockclusterInstance()
	res, err := bc.SendRequest(path, requestBody)

	log.Printf("G:TASK Fetching pods information: Response: %s", res)
}
