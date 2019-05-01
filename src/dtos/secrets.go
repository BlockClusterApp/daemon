package dtos

type WebappSecret struct {
	Metadata Metadata `json:"metadata"`
	Kind string `json:"kind"`
	ApiVersion string`json:"apiVersion"`
	Type string `json:"type"`
	Data struct {
		CaCert string `json:"ca.crt"`
		Namespace string `json:"namespace"`
		Token string `json:"token"`
	} `json:"data"`
}

type SecretList struct {
	Kind string `json:"kind"`
	Metadata Metadata `json:"metadata"`
	ApiVersion string `json:"apiVersion"`
	Items []WebappSecret `json:"items"`
}