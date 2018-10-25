package dtos

type Port struct {
	Name string `json:"name"`
	Protocol string `json:"protocol"`
	Port int32 `json:"port"`
	TargetPort int32 `json:"targetPort"`
	NodePort int32 `json:"nodePort"`
}
