package dtos

type Env struct {
	Name string `json:"name"`
	Value string `json:"value"`
	ValueFrom struct{
		SecretKeyRef struct{
			Name string `json:"name"`
			Key string `json:"key"`
		} `json:"secretKeyRef"`
		FieldRef struct{
			ApiVersion string `json:"apiVersion"`
			FieldPath string `json:"fieldPath"`
		} `json:"fieldRef"`
	} `json:"valueFrom"`
}

type Resources struct {
	Limits struct{
		Cpu string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"limits"`
	Requests struct{
		Cpu string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"requests"`
}

type VolumeMount struct{
	Name string `json:"name"`
	MountPath string `json:"mountPath"`
}

type Probe struct {
	Exec struct{
		Command []string `json:"command"`
	} `json:"exec"`
	HttpGet struct{
		Path string `json:"path"`
		Port int32 `json:"port"`
		Scheme string `json:"scheme"`
	} `json:"httpGet"`
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`
	TimeoutSeconds int32 `json:"timeoutSeconds"`
	PeriodSeconds int32 `json:"periodSeconds"`
	SuccessThreshold int32 `json:"successThreshold"`
	FailureThreshold int32 `json:"failureThreshold"`
}

type Container struct {
	Name string `json:"name"`
	Image string `json:"image"`
	WorkingDir string `json:"workingDir"`
	Ports []struct {
		Name string `json:"name"`
		ContainerPort int32 `json:"containerPort"`
		Protocol string `json:"protocol"`
	} `json:"ports"`
	Env []Env `json:"env"`
	Resources Resources `json:"resources"`
	VolumeMounts []VolumeMount `json:"volumeMounts"`
	LivelinessProbe Probe `json:"livelinessProbe"`
	ReadinessProbe Probe `json:"readinessProbe"`
	TerminationMessagePath string `json:"terminationMessagePath"`
	TerminationMessagePolicy string `json:"terminationMessagePolicy"`
	ImagePullPolicy string `json:"imagePullPolicy"`
}

type Volume struct {
	Name string `json:"name"`
	HostPath struct{
		Path string `json:"path"`
		Type string `json:"type"`
	} `json:"hostPath"`
	Secret struct{
		SecretName string `json:"secretName"`
		DefaultMode int32 `json:"defaultMode"`
	} `json:"secret"`
}

type ImagePullSecrets struct {
	Name string `json:"name"`
}

type Affinity struct {
	NodeAffinity struct{
		RequiredDuringSchedulingIgnoredDuringExecution struct{
			NodeSelectorTerms []struct{
				MatchExpressions []struct{
					Key string `json:"key"`
					Operator string `json:"operator"`
					Values []string `json:"values"`
				} `json:"matchExpressions"`
			} `json:"nodeSelectorTerms"`
		} `json:"requiredDuringSchedulingIgnoredDuringExecution"`
	} `json:"nodeAffinity"`
}

type Toleration struct {
	Key string `json:"key"`
	Operator string `json:"operator"`
	Effect string `json:"effect"`
	TolerationSeconds int32 `json:"tolerationSeconds"`
}

type ContainerStatus struct {
	Name string `json:"name"`
	State struct{
		Running struct{
			StartedAt string `json:"startedAt"`
		} `json:"running"`
	} `json:"struct"`
	LastState struct{
		Terminated struct{
			ExitCode int8 `json:"exitCode"`
			Reason string `json:"reason"`
			StartedAt string `json:"startedAt"`
			FinishedAt string `json:"finishedAt"`
			ContainerID string `json:"containerID"`
		} `json:"terminated"`
	} `json:"lastState"`
	Ready bool `json:"ready"`
	RestartCount int32 `json:"restartCount"`
	Image string `json:"image"`
	ImageID string `json:"imageID"`
	ContainerID string `json:"containerID"`
}

type PodSpec struct {
	// Pod Specific
	Volumes []Volume `json:"volumes"`
	Containers []Container `json:"containers"`
	RestartPolicy string `json:"restartPolicy"`
	TerminationGracePeriodSeconds int32 `json:"terminationGracePeriodSeconds"`
	DNSPolicy string `json:"dnsPolicy"`
	ServiceAccountName string `json:"serviceAccountName"`
	ServiceAccount string `json:"serviceAccount"`
	NodeName string `json:"nodeName"`
	ImagePullSecrets []ImagePullSecrets `json:"imagePullSecrets"`
	Affinity Affinity `json:"affinity"`
	SchedulerName string `json:"schedulerName"`
	Tolerations []Toleration `json:"tolerations"`
}