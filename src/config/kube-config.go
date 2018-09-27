package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// type AuthDetails struct {
// 	User     string `json:"user"`
// 	Password string `json:"pass"`
// }

// type LocationConfig struct {
// 	MasterApiHost string `json:"masterApiHost"`
// 	WorkerNodeIP  string `json:"workerNodeIP"`
// 	LocationCode  string `json:"locationCode`
// 	LocationName  string `json:"locationName"`
// 	Auth          AuthDetails
// }

// type NameSpaceConfig struct {
// 	Location LocationConfig
// }

// type ClusterConfig struct {
// 	NameSpace NameSpaceConfig
// }

// type BCKubeConfig struct {
// 	Clusters ClusterConfig `json:"clusters"`
// }

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
