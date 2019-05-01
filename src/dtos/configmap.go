package dtos


type BlockclusterConfigmap struct {
	Kind string `json:"kind"`
	Metadata Metadata `json:"metadata"`
	ApiVersion string `json:"apiVersion"`
	Data ServerClusterConfig `json:"data"`
}


type ServerClusterConfig struct {
	License       string `json:"licence.yaml"`
	ClusterConfig string `json:"cluster-config.json"`
	Config        string `json:"config.json"`
}