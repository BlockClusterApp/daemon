package helpers

import (
	"encoding/json"
	config2 "github.com/BlockClusterApp/daemon/src/config"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

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
	err := json.Unmarshal([]byte(config2.GetKubeConfig()), &config)

	if err != nil {
		GetLogger().Printf("Error parsing config for namespaces %s", err.Error())
		return []string{}
	}


	keys := reflect.ValueOf(config.Clusters).MapKeys()
	namespaces := make([]string, len(keys))


	for i:=0;i<len(keys);i++{
		namespaces[i] = keys[i].String()
	}

	GetLogger().Printf("Namespaces %s", namespaces)

	return namespaces
}