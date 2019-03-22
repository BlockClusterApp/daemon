package dtos

type KubeVersion struct {
	Major string `json:"major"`
	Minor string `json:"minor"`
	GitVersion string `json:"gitVersion"`
	GoVersion string `json:"goVersion"`
}