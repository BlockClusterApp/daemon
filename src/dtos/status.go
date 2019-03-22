package dtos

type KubeStatus struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata struct{}
	Status string `json:"status"`
	Message string `json:"message"`
	Reason string `json:"reason"`
	Details struct {
		Name string `json:"name"`
		Group string `json:"group"`
		Kind string `json:"kind"`
	} `json:"details"`
	Code int32 `json:"code"`
}