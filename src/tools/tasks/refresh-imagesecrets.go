package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/config"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"net/http"
)

func DeployWebApp(namespace string){
	helpers.GetLogger().Printf("Checking and deploying webapp in %s", namespace)
	helpers.CheckAndDeployWebapp(namespace)
}

func RefreshImagePullSecrets() {
	bc := helpers.GetBlockclusterInstance()


	var kubeConfig = dtos.ClusterConfig{}
	err := json.Unmarshal([]byte(config.GetKubeConfig()), &kubeConfig)

	if err != nil {
		helpers.GetLogger().Printf("Refresh Image Pull Secrets | Error unmarshalling kube config %s", err.Error())
		return
	}

	namespaces := helpers.GetNamespaces()

	helpers.GetLogger().Printf("Starting image pull secret refresh")
	authorizationToken := helpers.GetAuthorizationToken()


	if len(authorizationToken) == 0 {
		helpers.GetLogger().Printf("Error refreshing image pull secrets. No authorization token")
		return
	}

	secretName := "blockcluster-regsecret"


	for _, namespace := range namespaces {
		go func(namespace string) {
			var secretJSON = fmt.Sprintf(`{
    			"apiVersion": "v1",
    			"data": {
					".dockerconfigjson": "%s"
    			},
    			"kind": "Secret",
    			"metadata": {
        			"name": "%s",
					"namespace": "%s"
    			},
    			"type": "kubernetes.io/dockerconfigjson"
			}
			`, authorizationToken, secretName, namespace)


			locationCodes := helpers.GetLocationCodesOfEnv(kubeConfig.Clusters[namespace])

			for _, locationCode := range locationCodes {
				locationConfig := kubeConfig.Clusters[namespace][locationCode]


				path := fmt.Sprintf("%s/api/v1/namespaces/%s/secrets", locationConfig.MasterAPIHost, namespace)
				deletePath := fmt.Sprintf("%s/api/v1/namespaces/%s/secrets/%s", locationConfig.MasterAPIHost, namespace, secretName)

				deleteRequestOptions := helpers.ExternalKubeRequest{
					URL: deletePath,
					Auth: locationConfig.Auth,
					Method: http.MethodDelete,
				}

				_, err := helpers.MakeExternalKubeRequest(deleteRequestOptions)
				if err != nil {
					helpers.GetLogger().Printf("Error deleting image pull secrets %s", err.Error())
				}

				createRequestOptions := helpers.ExternalKubeRequest{
					URL: path,
					Auth: locationConfig.Auth,
					Method: http.MethodPost,
					Payload: secretJSON,
				}

				_, err = helpers.MakeExternalKubeRequest(createRequestOptions)

				if err != nil {
					helpers.GetLogger().Printf("Error refreshing image pull secrets %s", err.Error())
					continue
				}
			}


			helpers.GetLogger().Printf("Refreshed image pull secret in namespace %s", namespace)

			if bc.Metadata.ShouldDaemonDeployWebapp {
				DeployWebApp(namespace)
			}
		}(namespace)
	}

	helpers.GetLogger().Printf("Refreshed all image pull secrets")
}
