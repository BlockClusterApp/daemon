package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"net/http"
	"strings"
)



func FetchPod(selector string) *dtos.InfoResponse {
	var path = fmt.Sprintf("/api/v1/pods?labelSelector=%s", selector)
	response, err := MakeKubeRequest(http.MethodGet, path, nil)

	if err != nil {
		GetLogger().Printf("Error fetching pod details %s | %s",selector, err.Error())
		return nil
	}

	PodInfo := &dtos.InfoResponse{}
	err = json.Unmarshal([]byte(response), PodInfo)

	return PodInfo
}

func DeletePod(namespace string, podName string) bool {
	var path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName)
	_, err := MakeKubeRequest(http.MethodDelete, path, nil)

	if err != nil {
		GetLogger().Printf("Error deleting pod %s/%s | %s", namespace, podName, err.Error())
		return false
	}

	GetLogger().Printf("Deleted pod %s/%s", namespace, podName)

	return true
}

func FetchDeployment(selector string) *dtos.InfoResponse{
	var path = fmt.Sprintf("/apis/apps/v1beta2/deployments?labelSelector=%s", selector)
	response, err := MakeKubeRequest(http.MethodGet, path, nil)

	if err != nil {
		GetLogger().Printf("Error fetching deployment details %s | %s",selector, err.Error())
		return nil
	}

	DeployInfo := &dtos.InfoResponse{}
	err = json.Unmarshal([]byte(response), DeployInfo)

	return DeployInfo
}

func UpdateDeployment(deployInfo *dtos.InfoResponse) bool {
	var path = fmt.Sprintf("/apis/apps/v1beta2/namespaces/%s/deployment/%s", deployInfo.Metadata.Namespace, deployInfo.Metadata.Name)

	payload, err := json.Marshal(deployInfo)

	if err != nil {
		GetLogger().Printf("Error marshalling for deployment update %s/%s | %s", deployInfo.Metadata.Namespace, deployInfo.Metadata.Name, err.Error())
		return false
	}

	_, err2 := MakeKubeRequest(http.MethodPut, path, strings.NewReader(string(payload)))

	if err2 != nil {
		GetLogger().Printf("Error updating deployment %s/%s | %s", deployInfo.Metadata.Namespace, deployInfo.Metadata.Name, err.Error())
		return false
	}

	return true
}