package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
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
	GenerateName string `json:"generateName"`
	Namespace string `json:"namespace"`
	SelfLink string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
	Uid string `json:"uid"`
	Name string `json:"name"`
	CreationTimestamp string `json:"creationTimestamp"`
	Labels interface{} `json:"labels"`
	Annotations interface{} `json:"annotations"`
	OwnerReferences []struct{
		ApiVersion string `json:"apiVersion"`
		Kind string `json:"kind"`
		Name string `json:"name"`
		Uid string `json:"uid"`
		Controller bool `json:"controller"`
		BlockOwnerDeletion bool `json:"blockOwnerDeletion"`
	}
}

type InfoItems struct {
	Metadata Metadata `json:"metadata"`
	Spec NodeSpec `json:"spec"`
	Status NodeStatus `json:"status"`
}

type InfoResponse struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata Metadata `json:"metadata"`
	Items []InfoItems `json:"items"`
}

func FetchNodeInformation() {
	log := helpers.GetLogger()
	log.Println("G:TASK Fetching node information")
	nodeInfo, err := helpers.MakeKubeRequest(http.MethodGet, "/api/v1/nodes", nil)

	if err != nil {
		return
	}
	NodeMap := &InfoResponse{}
	err = json.Unmarshal([]byte(nodeInfo), NodeMap)

	log.Println(nodeInfo)

	if err != nil {
		log.Printf("Error parsing json while fetching node info %s", err.Error())
		return
	}

	requestBody := fmt.Sprintf(`{"info": "%s", "timestamp": "%d"}"`, nodeInfo, helpers.GetTimeInMillis())
	path := "/info/nodes"

	bc := helpers.GetBlockclusterInstance()
	res, err := bc.SendRequest(path, requestBody)

	log.Printf("G:TASK Fetching node information: Response: %s", res)
}