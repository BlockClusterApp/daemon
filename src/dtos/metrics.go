package dtos

type Usage struct {
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
}

type ContainerMetric struct {
	Name string `json:"name"`
	Usage Usage `json:"usage"`
}

type PodMetricItem struct {
	Metadata Metadata `json:"metadata"`
	Timestamp string `json:"timestamp"`
	Window string `json:"window"`
	Containers []ContainerMetric `json:"containers"`
}

type NodeMetricItem struct {
	Metadata Metadata `json:"metadata"`
	Timestamp string `json:"timestamp"`
	Window string `json:"window"`
	Usage Usage `json:"usage"`
}

type PodMetricResponse struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata Metadata `json:"metadata"`
	Items []PodMetricItem `json:"items"`
}

type NodeMetricResponse struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata Metadata `json:"metadata"`
	Items []NodeMetricItem `json:"items"`
}