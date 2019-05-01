package tools

import (
	"encoding/json"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/tools/tasks"
	"net/http"
	"strings"
)

func fetchBlockclusterConfigMap() (*dtos.BlockclusterConfigmap, error) {
	path := "/api/v1/namespaces/blockcluster/configmaps/blockcluster-config"

	config := &dtos.BlockclusterConfigmap{}

	res, err := helpers.MakeKubeRequest(http.MethodGet, path, nil)
	if err != nil {
		helpers.GetLogger().Printf("Error fetching blockcluster config : %s", err.Error())
		return config, err
	}

	err = json.Unmarshal([]byte(res), config)
	if err != nil {
		helpers.GetLogger().Printf("Error unmarshalling blockcluster configmap : %s", err.Error())
		return config, err
	}

	return config, nil
}

func deleteBlockclusterConfigmap() {
	path := "/api/v1/namespaces/blockcluster/configmaps/blockcluster-config"

	_, err := helpers.MakeKubeRequest(http.MethodDelete, path, nil)
	if err != nil {
		helpers.GetLogger().Printf("Error deleting blockcluster config : %s", err.Error())
	}
}

func fetchConfigFromServer() *dtos.ServerClusterConfig {
	bc := helpers.GetBlockclusterInstance()
	path := "/api/daemon/cluster-config"

	res, err := bc.SendGetRequest(path)

	result := &dtos.ServerClusterConfig{}

	if err != nil {
		helpers.GetLogger().Printf("Error fetching config from server: %s", err.Error())
		return result
	}

	err = json.Unmarshal([]byte(res), result)
	if err != nil {
		helpers.GetLogger().Printf("Error unmarshalling server cluster config : %s", err.Error())
		return result
	}

	return result
}

func createConfigMap(cm *dtos.BlockclusterConfigmap) {
	path := "/api/v1/namespaces/blockcluster/configmaps"

	body, err := json.Marshal(cm)
	if err != nil {
		helpers.GetLogger().Printf("Error converting configmap to json : %s", err.Error())
		return
	}

	_, err = helpers.MakeKubeRequest(http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		helpers.GetLogger().Printf("Error creating new configmap : %s", err.Error())
		return
	}

}

func Init() {
	tasks.ValidateLicence()

	tasks.SendWebappTokenToServer()

	_, err := fetchBlockclusterConfigMap()
	config := fetchConfigFromServer()

	kubeVersion := helpers.FetchLocalKubeVersion()

	var newConfig = dtos.BlockclusterConfigmap{
		ApiVersion: helpers.GetKubeAPIVersion(kubeVersion, "configmaps").APIVersion,
		Kind:       "ConfigMap",
		Metadata: dtos.Metadata{
			Name:      "blockcluster-config",
			Namespace: "blockcluster",
		},
		Data: dtos.ServerClusterConfig{
			License:       config.License,
			ClusterConfig: config.ClusterConfig,
			Config:        config.Config,
		},
	}

	if err != nil && (strings.Contains(err.Error(), "Not Found") || strings.Contains(err.Error(), "404")) {
		// Does not exists
	} else {
		deleteBlockclusterConfigmap()
	}

	createConfigMap(&newConfig)
}
