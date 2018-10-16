package dtos

type Strategy struct {
	Type string `json:"type"`
	RollingUpdate struct{
		MaxUnavailable string `json:"maxUnavailable"`
		MaxSurge string `json:"maxSurge"`
	} `json:"rollingUpdate"`
}

type Selector struct {
	MatchLabels struct{
		Name string `json:"name"`
	} `json:"matchLabels"`
}

type Template struct {
	Metadata Metadata `json:"metadata"`
	Spec PodSpec `json:"spec"`
	Status Status `json:"status"`
}