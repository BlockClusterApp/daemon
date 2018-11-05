package helpers

import (
	"encoding/json"
	config2 "github.com/BlockClusterApp/daemon/src/config"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
	"time"
)

const CURRENT_AGENT_VERSION = "1.0";

func UnmarshalJson(input []byte) (map[string]interface{}, error) {
	var result interface{}
	// unmarshal json to a map
	foomap := make(map[string]interface{})
	json.Unmarshal(input, &foomap)

	// create a mapstructure decoder
	var md mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Metadata: &md,
			Result:   result,
		})
	if err != nil {
		return nil, err
	}

	// decode the unmarshalled map into the given struct
	if err := decoder.Decode(foomap); err != nil {
		return nil, err
	}

	// copy and return unused fields
	unused := map[string]interface{}{}
	for _, k := range md.Unused {
		unused[k] = foomap[k]
	}
	return unused, nil
}

func GetTimeInMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetNamespaces() []string {
	var config = dtos.ClusterConfig{}
	err := json.Unmarshal([]byte(config2.GetRawKubeConfig()), &config)

	if err != nil {
		GetLogger().Printf("Error parsing config for namespaces %s", err.Error())
		return []string{}
	}

	keys := reflect.ValueOf(config.Clusters).MapKeys()
	namespaces := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		namespaces[i] = keys[i].String()
	}

	GetLogger().Printf("Namespaces %s", namespaces)

	return namespaces
}

func ReplaceWebAppConfig(fileContent string, webappConfig dtos.WebAppConfig, namespace string) string {
	replacer := strings.NewReplacer("%__NAMESPACE__%", namespace,
		"%__MONGO_URL__%", webappConfig.MongoConnectionURL,
		"%__REDIS_HOST__%", webappConfig.RedisHost,
		"%__REDIS_PORT__%", webappConfig.RedisPort,
		"%__IMAGE_URL__%", webappConfig.ImageRepository,
	)

	return replacer.Replace(fileContent)
}

func GetLocationCodesOfEnv(config map[string]*dtos.LocationConfig) []string {
	keys := reflect.ValueOf(config).MapKeys()
	locationCodes := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		locationCodes[i] = keys[i].String()
	}
	return locationCodes
}

type SimpleRepository struct {
	URL map[string]string `json:"url"`
}

type RepositoryConfig struct {
	Dynamo SimpleRepository `json:"dynamo"`
	Impulse SimpleRepository `json:"impulse"`
}

func GetRepositoryConfigForConfig() RepositoryConfig {
	var config = RepositoryConfig{}

	webAppConfig := config2.GetWebAppConfig()

	namespaces := reflect.ValueOf(webAppConfig.Dynamo).MapKeys()

	dynamoRepo := make(map[string]string, len(namespaces))
	impulseRepo := make(map[string]string, len(namespaces))

	for _,namespace := range namespaces {
		namespace := namespace.String()
		dynamoRepo[namespace] = webAppConfig.Dynamo[namespace]
		impulseRepo[namespace] = webAppConfig.Impulse[namespace]
	}

	config.Dynamo.URL = dynamoRepo
	config.Impulse.URL = impulseRepo

	return config
}