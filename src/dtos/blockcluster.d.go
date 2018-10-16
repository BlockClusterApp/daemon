package dtos

type LicenceValidationResponse struct {
	Success bool `json:"success"`
	Token string `json:"message"`
	Error string `json:"error"`
	ErrorCode int `json:"errorCode"`
	Metadata struct {
		BlockClusterAgentVersion string `json:"blockclusterAgentVersion"`
		WebAppVersion string `json:"webappVersion"`
		ShouldDaemonDeployWebapp bool `json:"shouldDaemonDeployWebapp"`
	} `json:"metadata"`
}

type AWSCredsResponse struct {
	ClientID string `json:"clientId"`
	AccessKeys struct{
		PolicyId string `json:"PolicyId"`
		AccessKeyId string `json:"AccessKeyId"`
		CreateDate string `json:"CreateDate"`
		SecretAccessKey string `json:"SecretAccessKey"`
		Status string `json:"Status"`
		UserName string `json:"UserName"`
	} `json:"accessKeys"`
}