package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/config"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/templates"
	"net/http"
	"strings"
)

func _deployHyperionScaler(locationConfig *dtos.LocationConfig, namespace string) {
	replacer := strings.NewReplacer("%__NAMESPACE__%", namespace,
		"%__K8S_HOST__%", locationConfig.MasterAPIHost,
		"%__K8S_USER__%", locationConfig.Auth.User,
		"%__K8S_PASS__%", locationConfig.Auth.Pass,
		"%__PROM_BASE_URI__%", "",
	)

	template := replacer.Replace(templates.GetHyperionScalerCronJobTemplate())

	helpers.GetLogger().Printf("Creating hyperion scaler deployment in %s/%s", namespace, locationConfig.LocationCode)
	url := fmt.Sprintf("%s/apis/batch/v1beta1/namespaces/%s/cronjobs", locationConfig.MasterAPIHost, namespace)

	go helpers.CreateResource(template, url, locationConfig.Auth)
}

func CheckHyperionScaler() {
	helpers.GetLogger().Printf("Starting hyperion scaler")
	var clusterConfig = dtos.ClusterConfig{}
	_ = json.Unmarshal([]byte(config.GetRawKubeConfig()), &clusterConfig)

	namespaces := helpers.GetNamespaces()

	for _, namespace := range namespaces {
		locationCodes := helpers.GetLocationCodesOfEnv(clusterConfig.Clusters[namespace])
		helpers.GetLogger().Printf("Starting hyperion scaler check for %s in %s", namespace, locationCodes)
		for _, locationCode := range locationCodes {
			go func(locationCode string, namespace string) {
				helpers.GetLogger().Printf("Checking hyperion scaler in %s/%s", namespace, locationCode)
				locationConfig := clusterConfig.Clusters[namespace][locationCode]

				var requestParams = helpers.ExternalKubeRequest{
					URL:     fmt.Sprintf("%s/apis/batch/v1beta1/namespaces/%s/cronjobs?fieldSelector=metadata.name%3Dhyperion-scaler", locationConfig.MasterAPIHost, namespace),
					Auth:    locationConfig.Auth,
					Payload: "",
					Method:  http.MethodGet,
				}

				resp, err := helpers.MakeExternalKubeRequest(requestParams)

				if err != nil {
					helpers.GetLogger().Printf("Error fetching cronjobs %s | %s", requestParams.URL, err.Error())
					return
				}

				var cronResponse = dtos.InfoResponse{}
				err = json.Unmarshal([]byte(resp), &cronResponse)

				if err != nil {
					helpers.GetLogger().Printf("Error unmarshalling cronjob response %s | %s", requestParams.URL, err.Error())
					return
				}

				if len(cronResponse.Items) <= 0 {
					go _deployHyperionScaler(locationConfig, namespace)
					return
				}

				helpers.GetLogger().Printf("Hyperion scaler exists in %s/%s", namespace, locationCode)
			}(locationCode, namespace)

		}

	}
}
