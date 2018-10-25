package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"net/http"
)

func FetchNodeInformation() {
	log := helpers.GetLogger()
	log.Println("G:TASK Fetching node information")
	nodeInfo, err := helpers.MakeKubeRequest(http.MethodGet, "/api/v1/nodes", nil)

	if err != nil {
		return
	}
	NodeMap := &dtos.InfoResponse{}
	err = json.Unmarshal([]byte(nodeInfo), NodeMap)

	log.Println(nodeInfo)

	if err != nil {
		log.Printf("Error parsing json while fetching node info %s", err.Error())
		return
	}

	requestBody := fmt.Sprintf(`{"info": %s, "timestamp": %d}`, nodeInfo, helpers.GetTimeInMillis())
	path := "/info/nodes"

	bc := helpers.GetBlockclusterInstance()
	res, err := bc.SendRequest(path, requestBody)

	if err != nil {
		log.Printf("Error sending node info %s", err.Error())
		return
	}

	log.Printf("G:TASK Fetching node information: Response: %s", res)
}
