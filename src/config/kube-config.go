package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetKubeConfig() string {

	fileAbsPath, _ := filepath.Abs("/conf.d/config.json")
	file, e := ioutil.ReadFile(fileAbsPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	// fmt.Printf("%s\n", string(file))
	return string(file)
}
