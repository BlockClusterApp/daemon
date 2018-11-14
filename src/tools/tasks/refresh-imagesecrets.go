package tasks

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"net/http"
	"strings"
)

func RefreshImagePullSecrets() {
	bc := helpers.GetBlockclusterInstance()

	namespaces := helpers.GetNamespaces()

	if bc.Metadata.ShouldDaemonDeployWebapp {
		for _, namespace := range namespaces {
			go func(namespace string) {
				helpers.GetLogger().Printf("Checking and deploying webapp in %s", namespace)
				helpers.CheckAndDeployWebapp(namespace)
			}(namespace)
		}

	}

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

			path := fmt.Sprintf("/api/v1/namespaces/%s/secrets", namespace)
			deletePath := fmt.Sprintf("/api/v1/namespaces/%s/secrets/%s", namespace, secretName)
			_, err := helpers.MakeKubeRequest(http.MethodDelete, deletePath, nil)
			if err != nil {
				helpers.GetLogger().Printf("Error deleting image pull secrets %s", err.Error())
			}

			_, err = helpers.MakeKubeRequest(http.MethodPost, path, strings.NewReader(secretJSON))

			if err != nil {
				helpers.GetLogger().Printf("Error refreshing image pull secrets %s", err.Error())
				return
			}
			helpers.GetLogger().Printf("Refreshed image pull secret in namespace %s", namespace)
		}(namespace)
	}

	helpers.GetLogger().Printf("Refreshed all image pull secrets")


}
