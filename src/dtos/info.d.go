package dtos

type Taint struct {
	Key    string `json:"key"`
	Effect string `json:"effect"`
}

type Spec struct {
	// Deployment Specific
	Replicas                int32    `json:"replicas"`
	RevisionHistoryLimit    int32    `json:"revisionHistoryLimit"`
	ProgressDeadlineSeconds int32    `json:"progressDeadlineSeconds"`
	Strategy                Strategy `json:"strategy"`
	Selector                Selector `json:"selector"`
	Template                Template `json:"template"`

	// Node Specific
	PodCIDR    string  `json:"podCIDR"`
	ExternalID string  `json:"externalID"`
	ProviderID string  `json:"providerID"`
	Taints     []Taint `json:"taints"`

	// Pod Specific
	Volumes                       []Volume           `json:"volumes"`
	Containers                    []Container        `json:"containers"`
	RestartPolicy                 string             `json:"restartPolicy"`
	TerminationGracePeriodSeconds int32              `json:"terminationGracePeriodSeconds"`
	DNSPolicy                     string             `json:"dnsPolicy"`
	ServiceAccountName            string             `json:"serviceAccountName"`
	ServiceAccount                string             `json:"serviceAccount"`
	NodeName                      string             `json:"nodeName"`
	ImagePullSecrets              []ImagePullSecrets `json:"imagePullSecrets"`
	Affinity                      Affinity           `json:"affinity"`
	SchedulerName                 string             `json:"schedulerName"`
	Tolerations                   []Toleration       `json:"tolerations"`

	// Service Specific
	Ports                 []Port `json:"ports"`
	ClusterIP             string `json:"clusterIP"`
	Type                  string `json:"type"`
	SessionAffinity       string `json:"sessionAffinity"`
	ExternalTrafficPolicy string `json:"externalTrafficPolicy"`
}

type Memory struct {
	Pods   string `json:"pods"`
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Condition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastHeartBeatTime  string `json:"lastHeartbeatTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastProbeTime      string `json:"lastProbeTime"`
}

type DaemonEndpoint struct {
	KubeletEndpoint struct {
		Port int64 `json:"Port"`
	}
}

type Address struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type NodeInfo struct {
	MachineID               string `json:"machineID"`
	SystemUUID              string `json:"systemUUID"`
	BootID                  string `json:"bootID"`
	KernelVersion           string `json:"kernelVersion"`
	OsImage                 string `json:"osImage"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KubeletVersion          string `json:"kubeletVersion"`
	KubeProxyVersion        string `json:"kubeProxyVersion"`
	OperatingSystem         string `json:"operatingSystem"`
	Architecture            string `json:"architecture"`
}

type Image struct {
	Names     []string `json:"names"`
	SizeBytes int64    `json:"sizeBytes"`
}

type Status struct {
	Capacity          Memory            `json:"capacity"`
	Allocatable       Memory            `json:"allocatable"`
	Conditions        []Condition       `json:"conditions"`
	Addresses         []Address         `json:"addresses"`
	DaemonEndpoints   DaemonEndpoint    `json:"daemonEndpoints"`
	Info              NodeInfo          `json:"info"`
	Images            []Image           `json:"images"`
	Phase             string            `json:"phase"`
	HostIP            string            `json:"hostIP"`
	PodIP             string            `json:"podIP"`
	StartTime         string            `json:"startTime"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses"`
	QosClass          string            `json:"qosClass"`
}

type Metadata struct {
	GenerateName      string                 `json:"generateName"`
	Namespace         string                 `json:"namespace"`
	SelfLink          string                 `json:"selfLink"`
	ResourceVersion   string                 `json:"resourceVersion"`
	Uid               string                 `json:"uid"`
	Name              string                 `json:"name"`
	Generation        int32                  `json:"generation"`
	CreationTimestamp string                 `json:"creationTimestamp"`
	Labels            interface{}            `json:"labels"`
	Annotations       map[string]interface{} `json:"annotations"`
	OwnerReferences   []struct {
		ApiVersion         string `json:"apiVersion"`
		Kind               string `json:"kind"`
		Name               string `json:"name"`
		Uid                string `json:"uid"`
		Controller         bool   `json:"controller"`
		BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
	}
}

type InfoItems struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

type InfoResponse struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   Metadata    `json:"metadata"`
	Items      []InfoItems `json:"items"`
}
