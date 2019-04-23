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

const CURRENT_AGENT_VERSION = "1.2";

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
		"%__RAZORPAY_ID__%", webappConfig.RazorPay.Id,
		"%__RAZORPAY_KEY__%", webappConfig.RazorPay.Key,
		"%__ROOT_URL__%", webappConfig.RootURL,
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

type SimplePrivatehiveRepository struct {
	URL map[string]dtos.PrivatehiveImages `json:"url"`
}

type RepositoryConfig struct {
	Dynamo SimpleRepository `json:"dynamo"`
	Impulse SimpleRepository `json:"impulse"`
	Privatehive SimplePrivatehiveRepository `json:"privatehive"`
}

func GetRepositoryConfigForConfig() (RepositoryConfig, dtos.WebAppConfigFile) {
	var config = RepositoryConfig{}

	webAppConfig := config2.GetWebAppConfig()

	namespaces := reflect.ValueOf(webAppConfig.Dynamo).MapKeys()

	dynamoRepo := make(map[string]string, len(namespaces))
	impulseRepo := make(map[string]string, len(namespaces))
	privatehiveRepo := make(map[string]dtos.PrivatehiveImages, len(namespaces))

	for _,namespace := range namespaces {
		namespace := namespace.String()
		dynamoRepo[namespace] = webAppConfig.Dynamo[namespace]
		impulseRepo[namespace] = webAppConfig.Impulse[namespace]
		privatehiveRepo[namespace] = dtos.PrivatehiveImages{
			Peer: webAppConfig.Privatehive[namespace].Peer,
			Orderer: webAppConfig.Privatehive[namespace].Orderer,
		}
	}

	config.Dynamo.URL = dynamoRepo
	config.Impulse.URL = impulseRepo
	config.Privatehive.URL = privatehiveRepo

	return config, webAppConfig
}

func TrimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func DoesArrayIncludeString(array []string, searchItem interface{}) bool {
	for _, item := range array {
		if item == searchItem {
			return true
		}
	}
	return false
}