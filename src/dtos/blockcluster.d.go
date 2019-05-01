package dtos

type LicenceMetadata struct {
	BlockClusterAgentVersion        string   `json:"blockclusterAgentVersion"`
	WebAppVersion                   string   `json:"webappVersion"`
	ShouldDaemonDeployWebapp        bool     `json:"shouldDaemonDeployWebapp"`
	ClientID                        string   `json:"clientId"`
	ActivatedFeatures               []string `json:"activatedFeatures"`
	ShouldWebAppRefreshAWSImageAuth bool     `json:"shouldWebAppRefreshAWSImageAuth"`
	WebAppMigration                 int32    `json:"webAppMigration"`
}

type LicenceValidationResponse struct {
	Success     bool            `json:"success"`
	Token       string          `json:"message"`
	Error       string          `json:"error"`
	ErrorCode   int             `json:"errorCode"`
	Metadata    LicenceMetadata `json:"metadata"`
	ClusterInfo []ClusterInfo     `json:"clusterInfo"`
}

type AWSCredsResponse struct {
	ClientID    string   `json:"clientId"`
	RegistryIds []string `json:"registryIds"`
	AccessKeys  struct {
		PolicyId        string `json:"PolicyId"`
		AccessKeyId     string `json:"AccessKeyId"`
		CreateDate      string `json:"CreateDate"`
		SecretAccessKey string `json:"SecretAccessKey"`
		Status          string `json:"Status"`
		UserName        string `json:"UserName"`
	} `json:"accessKeys"`
}

type ClusterInfo struct {
	MasterAPIHost    string         `json:"masterAPIHost"`
	WorkerNodeIP     string         `json:"workerNodeIP"`
	LocationCode     string         `json:"locationCode"`
	LocationName     string         `json:"locationName"`
	DynamoDomainName string         `json:"dynamoDomainName"`
	APIHost          string         `json:"apiHost"`
	Auth             Auth           `json:"auth"`
	Hyperion         HyperionConfig `json:"hyperion"`
}

type OperationType int

const (
	CLOUD_CONFIG OperationType = 1
	YAML_CONFIG OperationType = 2
)
