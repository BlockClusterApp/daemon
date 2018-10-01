package helpers

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LicenceConfig struct {
	Key string `yaml:"key"`
}

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

func GetLicenceKey() string {
	var log = GetLogger()
	var licence = LicenceConfig{}
	content := getLicenceFileContent()
	err := yaml.Unmarshal([]byte(content), &licence)
	if err != nil {
		log.Printf("Error reading licence key %s", err.Error())
	}
	return licence.Key
}
