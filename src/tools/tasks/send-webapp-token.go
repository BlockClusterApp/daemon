package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"net/http"
	"os"
	"strings"
)

func fetchSecrets() (*dtos.SecretList, error) {
	path := "/api/v1/namespaces/blockcluster/secrets?limit=500"
	res, err := helpers.MakeKubeRequest(http.MethodGet, path, nil)

	if err != nil {
		helpers.GetLogger().Printf("Error fetching SecretList %s", err.Error())
		return &dtos.SecretList{}, err
	}

	list := &dtos.SecretList{}

	err = json.Unmarshal([]byte(res), list)

	if err != nil {
		helpers.GetLogger().Printf("Error unmarshalling SecretList %s", err.Error())
		return &dtos.SecretList{}, err
	}

	return list, nil
}

func fetchWebappSecret(secretList *dtos.SecretList) dtos.WebappSecret {
	var webappSecret dtos.WebappSecret
	for _, secret := range secretList.Items {
		if strings.Contains(secret.Metadata.Name, "blockcluster-webapp-token") {
			webappSecret = secret
			break
		}
	}

	return webappSecret
}

func sendTokenToServer(token string) {
	path := "/api/daemon/cluster-token"
	blockcluster := helpers.GetBlockclusterInstance()

	body := fmt.Sprintf(`{"token": "%s", "identifier": "%s"}`, token, os.Getenv("CLUSTER_IDENTIFIER"))
	res, err := blockcluster.SendRequest(path, body)
	if err !=nil {
		helpers.GetLogger().Printf("Error sending token : %s", err.Error())
		return
	}
	helpers.GetLogger().Printf("Sent token to server : %s", res)
}

func SendWebappTokenToServer() {
	secretList, err := fetchSecrets()
	if err != nil {
		return
	}

	webappSecret := fetchWebappSecret(secretList)

	if len(webappSecret.Data.Token) <= 0 {
		helpers.GetLogger().Printf("Error fetching webapp secret from secret list")
		return
	}

	sendTokenToServer(webappSecret.Data.Token)
	helpers.GetLogger().Printf("Sent token")
}