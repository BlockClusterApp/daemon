package tasks

import (
	"encoding/json"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"log"
	"net/http"
)

type NodeTaint struct {
	Key string `json:"key"`
	Effect string `json:"effect"`
}

type NodeSpec struct{
	PodCIDR string `json:"podCIDR"`
	ExternalID string `json:"externalID"`
	ProviderID string `json:"providerID"`
	Taints []NodeTaint `json:"taints"`
}

type NodeMemory struct {
	Pods string `json:"pods"`
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
}

type NodeCondition struct {
	Type string `json:"type"`
	Status string `json:"status"`
	LastHeartBeatTime string `json:"lastHeartbeatTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
	Reason string `json:"reason"`
	Message string `json:"message"`
}

type NodeDaemonEndpoint struct {
	KubeletEndpoint struct{
		Port int64 `json:"Port"`
	}
}

type NodeAddress struct {
	Type string `json:"type"`
	Address string `json:"address"`
}

type NodeInfo struct {
	MachineID string `json:"machineID"`
	SystemUUID string `json:"systemUUID"`
	BootID string `json:"bootID"`
	KernelVersion string `json:"kernelVersion"`
	OsImage string `json:"osImage"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KubeletVersion string `json:"kubeletVersion"`
	KubeProxyVersion string `json:"kubeProxyVersion"`
	OperatingSystem string `json:"operatingSystem"`
	Architecture string `json:"architecture"`
}

type NodeImage struct {
	Names []string `json:"names"`
	SizeBytes int64 `json:"sizeBytes"`
}

type NodeStatus struct {
	Capacity NodeMemory `json:"capacity"`
	Allocatable NodeMemory `json:"allocatable"`
	Conditions []NodeCondition `json:"conditions"`
	Addresses []NodeAddress `json:"addresses"`
	DaemonEndpoints NodeDaemonEndpoint `json:"daemonEndpoints"`
	NodeInfo NodeInfo `json:"nodeInfo"`
	Images []NodeImage `json:"images"`
}


type Metadata struct {
	SelfLink string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
	Uid string `json:"uid"`
	Name string `json:"name"`
	CreationTimestamp string `json:"creationTimestamp"`
	Labels interface{} `json:"labels"`
	Annotations interface{} `json:"annotations"`
}

type InfoItems struct {
	Metadata Metadata `json:"metadata"`
	Spec NodeSpec `json:"spec"`
	Status NodeStatus `json:"status"`
}

type NodeInfoResponse struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata Metadata `json:"metadata"`
	Items []InfoItems `json:"items"`
}

func FetchNodeInformation() {
	log.Println("GOCRON:TASK Fetching node information")
	nodeInfo, err := helpers.MakeRequest(http.MethodGet, "/api/v1/nodes", nil)

	if err != nil {
		return
	}
	NodeMap := &NodeInfoResponse{}
	err = json.Unmarshal([]byte(nodeInfo), NodeMap)

	if err != nil {
		log.Printf("Error parsing json while fetching node info %s", err.Error())
		return
	}

	log.Printf("Node Info %s", NodeMap)
}