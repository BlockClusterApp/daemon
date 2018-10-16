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

	var secretJSON = fmt.Sprintf(`
		{
    		"apiVersion": "v1",
    		"data": {
				".dockerconfigjson": "%s"
    		},
    		"kind": "Secret",
    		"metadata": {
        		"name": "regsecret",
        		"namespace": "default"
    		},
    		"type": "kubernetes.io/dockerconfigjson"
		}
	`, authorizationToken)

	path := fmt.Sprintf("/api/v1/namespaces/%s/secrets", "default")
	_, err := helpers.MakeKubeRequest(http.MethodPut, path, strings.NewReader(secretJSON))

	if err != nil {
		helpers.GetLogger().Printf("Error refreshing image pull secrets %s", err.Error())
		return
	}

	helpers.GetLogger().Printf("Refreshed image pull secrets")
}