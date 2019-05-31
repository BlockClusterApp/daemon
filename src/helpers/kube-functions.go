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

func FetchKubeVersion(locationConfig dtos.LocationConfig) *dtos.KubeVersion {
	var requestParams = ExternalKubeRequest{
		URL:     fmt.Sprintf("%s/version", locationConfig.MasterAPIHost),
		Auth:    locationConfig.Auth,
		Payload: "",
		Method:  http.MethodGet,
	}

	resp, err := MakeExternalKubeRequest(requestParams)

	if err != nil {
		GetLogger().Printf("Error fetching kube version details %s | %s", locationConfig.MasterAPIHost, err.Error())
		return nil
	}

	versionInfo := &dtos.KubeVersion{}
	err = json.Unmarshal([]byte(resp), versionInfo)

	return versionInfo
}

func FetchLocalKubeVersion() *dtos.KubeVersion {
	resp, err := MakeKubeRequest(http.MethodGet, "/version", nil)
	if err != nil {
		GetLogger().Printf("Error fetching local kube version details | %s", err.Error())
		return nil
	}

	versionInfo := &dtos.KubeVersion{}
	err = json.Unmarshal([]byte(resp), versionInfo)

	return versionInfo
}

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

func CreateResource(config string, url string, auth dtos.Auth) {
	params := ExternalKubeRequest{
		URL:     url,
		Auth:    auth,
		Payload: config,
		Method:  http.MethodPost,
	}

	resp, err := MakeExternalKubeRequest(params)

	if err != nil {
		GetLogger().Printf("Error creating resource %s | %s", params, err.Error())
		return
	}

	GetLogger().Printf("Successfully created resource %s \n %s", params, resp)
	return
}

func _checkAndDeployWebapp(namespace string, locationConfig dtos.LocationConfig, webAppConfig dtos.WebAppConfig) {
	GetLogger().Printf("Checking and deploying webapp for %s/%s", namespace, locationConfig.LocationCode)

	path := fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments?labelSelector=%s", namespace, "app%3Dblockcluster-app")
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
	var deployment = dtos.InfoResponse{}
	err = json.Unmarshal([]byte(response), &deployment)

	if err != nil {
		GetLogger().Printf("Error parsing deployment response to struct %s | %s", url, err.Error())
		return
	}

	url = fmt.Sprintf("%s/api/v1/namespaces/%s/services?fieldSelector=metadata.name%%blockcluster-svc", locationConfig.MasterAPIHost, namespace)
	params.URL = url
	response, err = MakeExternalKubeRequest(params)
	var service = dtos.InfoResponse{}
	err = json.Unmarshal([]byte(response), &service)

	url = fmt.Sprintf("%s/apis/autoscaling/v1/namespaces/%s/horizontalpodautoscalers?fieldSelector=metadata.name%%blockcluster-hpa", locationConfig.MasterAPIHost, namespace)
	params.URL = url
	response, err = MakeExternalKubeRequest(params)
	var hpa = dtos.InfoResponse{}
	err = json.Unmarshal([]byte(response), &hpa)

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

	deployURL := fmt.Sprintf("%s/apis/apps/v1beta1/namespaces/%s/deployments", locationConfig.MasterAPIHost, namespace)
	serviceURL := fmt.Sprintf("%s/api/v1/namespaces/%s/services", locationConfig.MasterAPIHost, namespace)
	hpaURL := fmt.Sprintf("%s/apis/autoscaling/v1/namespaces/%s/horizontalpodautoscalers", locationConfig.MasterAPIHost, namespace)

	go CreateResource(deploymentConfig, deployURL, locationConfig.Auth)
	go CreateResource(serviceConfig, serviceURL, locationConfig.Auth)
	go CreateResource(hpaConfig, hpaURL, locationConfig.Auth)

}

func CheckAndDeployWebapp(namespace string) {

	var kubeConfig = dtos.ClusterConfig{}
	err := json.Unmarshal([]byte(config2.GetKubeConfig()), &kubeConfig)

	if err != nil {
		GetLogger().Printf("Check and Deploy webapp | Error unmarshalling kube config %s", err.Error())
		return
	}

	config := config2.GetWebAppConfig()

	namespaces := reflect.ValueOf(kubeConfig.Clusters).MapKeys()

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
		RedisHost:          config.Redis[namespace].Host,
		RedisPort:          config.Redis[namespace].Port,
		ImageRepository:    config.WebApp[namespace],
		RootURL:            config.RootUrl[namespace],
	}

	for _, locationCode := range locationCodes {
		locationConfig := kubeConfig.Clusters[namespace][locationCode]
		go _checkAndDeployWebapp(namespace, *locationConfig, _webAppConfig)
	}

}

func _checkAndDeployHyperion(namespace string, locationConfig dtos.LocationConfig) {
	GetLogger().Printf("Checking and deploying hyperion for %s/%s", locationConfig.LocationCode, namespace)

	kubeVersion := FetchKubeVersion(locationConfig)
	if kubeVersion == nil {
		kubeVersion = &dtos.KubeVersion{
			Major:      "1",
			Minor:      "9",
			GitVersion: "v1.9.8",
		}
	}
	statefulSetMapping := GetKubeAPIVersion(kubeVersion, "statefulsets")
	checkHyperionPath := fmt.Sprintf("%s/%s/namespaces/%s/statefulsets/%s", locationConfig.MasterAPIHost, statefulSetMapping.Path, namespace, "hyperion")
	res, err := MakeExternalKubeRequest(ExternalKubeRequest{
		URL:     checkHyperionPath,
		Auth:    locationConfig.Auth,
		Payload: "",
		Method:  http.MethodGet,
	})

	statefulset := &dtos.InfoResponse{}
	err = json.Unmarshal([]byte(res), statefulset)

	if statefulset.Metadata.Name == "hyperion" {
		// Already exists
		GetLogger().Printf("Hyperion statefulset already exists in %s/%s", locationConfig.LocationCode, namespace)
		return
	}

	// Create cluster role
	clusterRoleMapping := GetKubeAPIVersion(kubeVersion, "clusterroles")
	path := fmt.Sprintf("%s/%s/clusterroles", locationConfig.MasterAPIHost, clusterRoleMapping.Path)
	req := ExternalKubeRequest{
		Auth:    locationConfig.Auth,
		Method:  http.MethodPost,
		URL:     path,
		Payload: templates.GetHyperionClusterRole(clusterRoleMapping.APIVersion),
	}

	_, err = MakeExternalKubeRequest(req)
	if err != nil {
		GetLogger().Printf("Error creating cluster role for hyperion: %s", err.Error())
	} else {
		GetLogger().Printf("Created hyperion cluster role in %s/%s", locationConfig.LocationCode, namespace)
	}

	// Create cluster role binding
	clusterRoleBindingMapping := GetKubeAPIVersion(kubeVersion, "clusterrolebindings")
	path = fmt.Sprintf("%s/%s/clusterrolebindings", locationConfig.MasterAPIHost, clusterRoleBindingMapping.Path)
	req = ExternalKubeRequest{
		Auth:    locationConfig.Auth,
		Method:  http.MethodPost,
		URL:     path,
		Payload: templates.GetHyperionClusterRoleBinding(clusterRoleBindingMapping.APIVersion, namespace),
	}
	_, err = MakeExternalKubeRequest(req)
	if err != nil {
		GetLogger().Printf("Error creating cluster role bindings for hyperion: %s", err.Error())
	} else {
		GetLogger().Printf("Created hyperion cluster role binding in %s/%s", locationConfig.LocationCode, namespace)
	}

	// Create stateful set

	path = fmt.Sprintf("%s/%s/namespaces/%s/statefulsets", locationConfig.MasterAPIHost, statefulSetMapping.Path, namespace)
	req = ExternalKubeRequest{
		Auth:    locationConfig.Auth,
		Method:  http.MethodPost,
		URL:     path,
		Payload: templates.GetHyperionStatefulSet(100, statefulSetMapping.APIVersion),
	}
	_, err = MakeExternalKubeRequest(req)
	if err != nil {
		GetLogger().Printf("Error creating hyperion statefulset: %s", err.Error())
	} else {
		GetLogger().Printf("Created hyperion stateful set in %s/%s", locationConfig.LocationCode, namespace)
	}

	// Create stateful set
	serviceMapping := GetKubeAPIVersion(kubeVersion, "services")
	path = fmt.Sprintf("%s/%s/namespaces/%s/services", locationConfig.MasterAPIHost, serviceMapping.Path, namespace)
	req = ExternalKubeRequest{
		Auth:    locationConfig.Auth,
		Method:  http.MethodPost,
		URL:     path,
		Payload: templates.GetHyperionService(serviceMapping.APIVersion),
	}
	_, err = MakeExternalKubeRequest(req)
	if err != nil {
		GetLogger().Printf("Error creating hyperion service: %s", err.Error())
	} else {
		GetLogger().Printf("Created hyperion service in %s/%s", locationConfig.LocationCode, namespace)
	}

	GetLogger().Printf("Created hyperion Statefulset")
}

func CheckAndDeployHyperion(namespace string) {

	if !DoesLicenceIncludeFeature("Hyperion") {
		GetLogger().Printf("Hyperion not in licence. Skipping deployment")
		return
	}

	var kubeConfig = dtos.ClusterConfig{}
	err := json.Unmarshal([]byte(config2.GetKubeConfig()), &kubeConfig)

	if err != nil {
		GetLogger().Printf("Check and Deploy Hyperion | Error unmarshalling kube config %s", err.Error())
		return
	}

	locationCodes := GetLocationCodesOfEnv(kubeConfig.Clusters[namespace])

	for _, locationCode := range locationCodes {
		locationConfig := kubeConfig.Clusters[namespace][locationCode]
		if locationConfig.Hyperion.IpfsPort != "" {
			go _checkAndDeployHyperion(namespace, *locationConfig)
		} else {
			GetLogger().Printf("Hyperion not in location %s", locationCode)
		}
	}
}
