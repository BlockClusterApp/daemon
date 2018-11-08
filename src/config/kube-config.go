package config

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/funcs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetKubeConfig() string {

	filePath, _ := filepath.Abs("/tmp/current-config.json")
	isFileEncrypted := true

	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		filePath, _ = filepath.Abs("/conf.d/cluster-config.json")
		isFileEncrypted = false
	}

	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	if isFileEncrypted {
		return funcs.DecryptString(string(file))
	}

	return string(file)
}

func GetRawKubeConfig() string {

	filePath, _ := filepath.Abs("/conf.d/cluster-config.json")
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return string(file)
}
