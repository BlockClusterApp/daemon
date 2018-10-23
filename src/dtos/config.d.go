package dtos

type LocationConfig struct {
	MasterAPIHost    string `json:"masterApiHost"`
	WorkerNodeIP     string `json:"workerNodeIP"`
	LocationCode     string `json:"locationCode"`
	LocationName     string `json:"locationName"`
	DynamoDomainName string `json:"dynamoDomainName"`
	APIHost          string `json:"apiHost"`
	Auth             struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"auth"`
	Hyperion struct {
		IpfsPort        string `json:"ipfsPort"`
		IpfsClusterPort string `json:"ipfsClusterPort"`
	} `json:"hyperion"`
}

type ClusterConfig struct {
	Clusters map[string]map[string]LocationConfig `json:"clusters"`
}