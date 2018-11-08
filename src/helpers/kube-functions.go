package helpers

import (
	"encoding/json"
	"fmt"
	config2 "github.com/BlockClusterApp/daemon/src/config"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/templates"
	"github.com/getsentry/raven-go"
	"net/http"
	"reflect"
	"strings"
)

func FetchPod(selector string) *dtos.InfoResponse {
	var path = fmt.Sprintf("/api/v1/pods?labelSelector=%s", selector)
	response, err := MakeKubeRequest(http.MethodGet, path, nil)

	if err != nil {
		bc := GetBlockclusterInstance()
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
		})
		GetLogger().Printf("Error fetching pod details %s | %s", selector, err.Error())
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
		bc := GetBlockclusterInstance()
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
		})
		GetLogger().Printf("Error deleting pod %s/%s | %s", namespace, podName, err.Error())
		return false
	}

	GetLogger().Printf("Deleted pod %s/%s", namespace, podName)

	return true
}

func fetchDeployment(path string) *dtos.InfoResponse {
	response, err := MakeKubeRequest(http.MethodGet, path, nil)

	if err != nil {
		GetLogger().Printf("Error fetching deployment details %s | %s", path, err.Error())
		return nil
	}

	DeployInfo := &dtos.InfoResponse{}
	err = json.Unmarshal([]byte(response), DeployInfo)

	return DeployInfo
}

func FetchDeployment(selector string) *dtos.InfoResponse {
	var path = fmt.Sprintf("/apis/apps/v1beta2/deployments?labelSelector=%s", selector)
	return fetchDeployment(path)
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
		bc := GetBlockclusterInstance()
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
		})
		GetLogger().Printf("Error updating deployment %s/%s | %s", deployInfo.Metadata.Namespace, deployInfo.Metadata.Name, err.Error())
		return false
	}

	return true
}

func createResource(config string, url string, auth dtos.Auth) {
	params := ExternalKubeRequest{
		URL:     url,
		Auth:    auth,
		Payload: config,
		Method:  http.MethodPost,
	}

	resp, err := MakeExternalKubeRequest(params)

	if err != nil {
		GetLogger().Printf("Error creating resource %s", params)
		return
	}

	GetLogger().Printf("Successfully created resource %s \n %s", params, resp)
	return
}

func _checkAndDeployWebapp(namespace string, locationConfig dtos.LocationConfig, webAppConfig dtos.WebAppConfig) {
	GetLogger().Printf("Checking and deploying webapp for %s/%s", namespace, locationConfig.LocationCode)

	path := fmt.Sprintf("/apis/apps/v1beta2/namespaces/%s/deployments?labelSelector=%s", namespace, "app%3Dblockcluster-app")
	url := fmt.Sprintf("%s%s", locationConfig.MasterAPIHost, path)
	params := ExternalKubeRequest{
		URL:     url,
		Auth:    locationConfig.Auth,
		Payload: "",
		Method:  http.MethodGet,
	}

	response, err := MakeExternalKubeRequest(params)
	if err != nil {
		return
	}
	var deployment dtos.InfoResponse
	err = json.Unmarshal([]byte(response), response)

	if err != nil {
		GetLogger().Printf("Error parsing deployment response to struct %s | %s", url, err.Error())
		return
	}

	url = fmt.Sprintf("%s/api/v1/namespaces/%s/services?fieldSelector=metadata.name%%blockcluster-svc", locationConfig.MasterAPIHost, namespace)
	params.URL = url
	response, err = MakeExternalKubeRequest(params)
	var service dtos.InfoResponse
	err = json.Unmarshal([]byte(response), service)

	url = fmt.Sprintf("%s/apis/autoscaling/v1/namespaces/%s/horizontalpodautoscalers?fieldSelector=metadata.name%%blockcluster-hpa", locationConfig.MasterAPIHost, namespace)
	params.URL = url
	response, err = MakeExternalKubeRequest(params)
	var hpa dtos.InfoResponse
	err = json.Unmarshal([]byte(response), hpa)

	if err != nil {
		GetLogger().Printf("Error parsing deployment response to struct %s | %s", url, err.Error())
		return
	}

	if len(deployment.Items) > 0 && len(service.Items) > 0 && len(deployment.Items) > 0 {
		// Deployment, Service and HPA is alredy present. No need to install it
		return
	}

	GetLogger().Printf("Deploying webapp for %s/%s", namespace, locationConfig.LocationCode)

	deploymentConfig := ReplaceWebAppConfig(templates.GetWebappDeploymentTemplate(), webAppConfig, namespace)
	serviceConfig := ReplaceWebAppConfig(templates.GetWebappServiceTemplate(), webAppConfig, namespace)
	hpaConfig := ReplaceWebAppConfig(templates.GetWebappHPATemplate(), webAppConfig, namespace)

	deployURL := fmt.Sprintf("%s/apis/apps/v1beta2/namespaces/%s/deployments", locationConfig.MasterAPIHost, namespace)
	serviceURL := fmt.Sprintf("%s/api/v1/namespaces/%s/services", locationConfig.MasterAPIHost, namespace)
	hpaURL := fmt.Sprintf("%s/apis/autoscaling/v1/namespaces/%s/horizontalpodautoscalers", locationConfig.MasterAPIHost, namespace)

	go createResource(deploymentConfig, deployURL, locationConfig.Auth)
	go createResource(serviceConfig, serviceURL, locationConfig.Auth)
	go createResource(hpaConfig, hpaURL, locationConfig.Auth)

}

func CheckAndDeployWebapp(namespace string) {

	var kubeConfig = dtos.ClusterConfig{}
	err := json.Unmarshal([]byte(config2.GetKubeConfig()), &kubeConfig)

	if err != nil {
		GetLogger().Printf("Check and Deploy webapp | Error unmarshalling kube config %s", err.Error())
		return
	}

	config := config2.GetWebAppConfig()

	namespaces := reflect.ValueOf(config).MapKeys()

	if len(namespaces) == 0 {
		GetLogger().Printf("Tried to deploy webapp to %s but webapp config has no namespace | %s", namespace, config)
		return
	}

	doesIncludeNamespace := false

	for _, name := range namespaces {
		if name.String() == namespace {
			doesIncludeNamespace = true
			break
		}
	}

	if !doesIncludeNamespace {
		GetLogger().Printf("Checking for %s but found it wasn't found in %s", namespace, namespaces)
		return
	}


	locationCodes := GetLocationCodesOfEnv(kubeConfig.Clusters[namespace])

	var _webAppConfig = dtos.WebAppConfig{
		MongoConnectionURL: config.MongoURL[namespace],
		RedisHost: config.Redis[namespace].Host,
		RedisPort: config.Redis[namespace].Port,
		ImageRepository: config.WebApp[namespace],
	}

	for _, locationCode := range locationCodes {
		locationConfig := kubeConfig.Clusters[namespace][locationCode]
		go _checkAndDeployWebapp(namespace, *locationConfig, _webAppConfig)
	}

}
