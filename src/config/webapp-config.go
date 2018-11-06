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

func GetWebAppConfig() dtos.WebAppConfigFile {

	var webAppConfig = dtos.WebAppConfigFile{}
	filePath, _ := filepath.Abs("/conf.d/config.json")
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}


	err := json.Unmarshal(file, &webAppConfig)

	if err != nil {
		log.Printf("Error parsing webapp config %s", err.Error())
	}


	log.Printf("Config file %s", webAppConfig.Dynamo)

	return webAppConfig
}
