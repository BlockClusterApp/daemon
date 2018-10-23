package tasks

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"net/http"
	"strings"
)

func RefreshImagePullSecrets() {
	helpers.GetLogger().Printf("Starting image pull secret refresh")
	authorizationToken := helpers.GetAuthorizationToken()

	if len(authorizationToken) == 0 {
		helpers.GetLogger().Printf("Error refreshing image pull secrets. No authorization token")
		return
	}

	secretName := "blockcluster-regsecret"
	namespace := "default"

	var secretJSON = fmt.Sprintf(`{
    		"apiVersion": "v1",
    		"data": {
				".dockerconfigjson": "%s"
    		},
    		"kind": "Secret",
    		"metadata": {
        		"name": "%s",
        		"namespace": "default"
    		},
    		"type": "kubernetes.io/dockerconfigjson"
		}
	`, authorizationToken, secretName)

	helpers.GetLogger().Printf("Secret %s", secretJSON)
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

	helpers.GetLogger().Printf("Refreshed image pull secrets")
}