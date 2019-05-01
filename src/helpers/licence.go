package helpers

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LicenceConfig struct {
	Key string `yaml:"key"`
}

var Licence LicenceConfig

func getLicenceFileContent() string {
	fileAbsPath, _ := filepath.Abs("/conf.d/licence.yaml")
	file, e := ioutil.ReadFile(fileAbsPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	// fmt.Printf("%s\n", string(file))
	return string(file)
}

func getLicenceKey() LicenceConfig {
	bc := GetBlockclusterInstance()
	var licence = LicenceConfig{}
	if bc.IsInInitMode{
		licence = LicenceConfig{
			Key: os.Getenv("LICENSE_KEY"),
		}
	} else {
		content := getLicenceFileContent()
		err := yaml.Unmarshal([]byte(content), &licence)
		if err != nil {
			raven.CaptureError(err, map[string]string{})
			GetLogger().Printf("Error reading licence key %s", err.Error())
		}
	}

	return licence
}

func UpdateLicence() {
	licence := getLicenceKey()
	Licence.Key = licence.Key
}

func GetLicence() LicenceConfig {
	return Licence
}

func DoesLicenceIncludeFeature(feature string) bool {
	features := GetBlockclusterInstance().Metadata.ActivatedFeatures

	return DoesArrayIncludeString(features, feature)
}
