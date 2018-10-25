package tasks

import (
	"encoding/json"
	"fmt"
	config2 "github.com/BlockClusterApp/daemon/src/config"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/tidwall/sjson"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

func UpdateHyperionPorts(){

	namespaces := helpers.GetNamespaces()

	var config = dtos.ClusterConfig{}
	err := json.Unmarshal([]byte(config2.GetRawKubeConfig()), &config)

	var newConfig = config2.GetRawKubeConfig()

	if err != nil {
		helpers.GetLogger().Printf("Error parsing config for namespaces %s", err.Error())
		return
	}


	// Need to be synchronous else file writing may contain inconsistent values. Or if you have any better way of
	for _, namespace := range namespaces {

		keys := reflect.ValueOf(config.Clusters[namespace]).MapKeys()
		locationCodes := make([]string, len(keys))

		for i:=0;i<len(keys);i++{
			locationCodes[i] = keys[i].String()
		}

		for _, locationCode := range locationCodes {
			locationConfig := config.Clusters[namespace][locationCode]
			var requestParams = helpers.ExternalKubeRequest{
				URL: fmt.Sprintf("%s/api/v1/namespaces/%s/services?fieldSelector=metadata.name%%3Dhyperion", locationConfig.MasterAPIHost, namespace),
				Auth: locationConfig.Auth,
				Payload: "",
				Method: http.MethodGet,
			}

			helpers.GetLogger().Printf("Fetching hyperion service details for (%s,%s) | %s", namespace, locationCode, requestParams.URL)

			resp, err := helpers.MakeExternalKubeRequest(requestParams)

			if err != nil {
				continue
			}

			var serviceResponse = dtos.InfoResponse{}
			err = json.Unmarshal([]byte(resp), &serviceResponse)

			if err != nil {
				helpers.GetLogger().Printf("Update Hyperion: Error unmarshalling service response %s", err.Error())
				continue
			}

			if len(serviceResponse.Items) == 0 {
				continue
			}

			hyperionService := serviceResponse.Items[0]

			for _,port := range hyperionService.Spec.Ports {
				if port.Name == "cluster-gateway" {
					value,_ := sjson.Set(newConfig, fmt.Sprintf("clusters.%s.%s.hyperion.ipfsPort", namespace, locationCode), strconv.Itoa(int(port.NodePort)))
					newConfig = value
				}
				if port.Name == "cluster-api" {
					value,_ := sjson.Set(newConfig, fmt.Sprintf("clusters.%s.%s.hyperion.ipfsClusterPort", namespace, locationCode), strconv.Itoa(int(port.NodePort)))
					newConfig = value
				}
			}
		}
	}


	outputFilePath := "/conf.d/current-config.json"

	 _, err = os.Stat(outputFilePath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(outputFilePath)
		if err != nil {
			helpers.GetLogger().Printf("Error creating file %s  | %s", outputFilePath, err.Error())
			return
		}
		defer file.Close()
	}


	file, err := os.OpenFile(outputFilePath, os.O_RDWR, 0644)

	if err != nil {
		helpers.GetLogger().Printf("Error opening file to write %s | %s ", outputFilePath, err.Error())
		return
	}

	defer file.Close()

	_, err = file.Write([]byte(newConfig))

	if err != nil {
		helpers.GetLogger().Printf("Error writing to file %s | %s", outputFilePath, err.Error())
		file.Close()
		return
	}

	err = file.Sync()
	if err != nil {
		helpers.GetLogger().Printf("Error syncing file %s | %s", outputFilePath, err.Error())
		file.Close()
		return
	}

	helpers.GetLogger().Printf("Successfully written new config")
}