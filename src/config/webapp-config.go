package config

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetWebAppConfig() map[string]dtos.WebAppConfigFile {

	var webAppConfig = map[string]dtos.WebAppConfigFile{}
	filePath, _ := filepath.Abs("/conf.d/config.json")
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, webAppConfig)

	if err != nil {
		log.Printf("Error parsing webapp config %s", err.Error())
	}

	return webAppConfig
}
