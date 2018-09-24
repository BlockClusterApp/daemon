package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/creamdog/gonfig"
)

func GetKubeConfig() gonfig.Gonfig {
	fileAbsPath, _ := filepath.Abs("./src/config-files/config.json")

	log.Println("Reading config from file", fileAbsPath)

	file, _ := os.Open(fileAbsPath)

	defer file.Close()

	config, _ := gonfig.FromJson(file)

	log.Println("Got config ", config)

	return config
}
