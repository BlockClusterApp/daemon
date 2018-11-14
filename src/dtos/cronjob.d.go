package dtos

type JobTemplate struct {
	Metadata struct {
		CreationTimestamp string `json:"creationTimestamp"`
	}
	Containers                    []Container        `json:"containers"`
	RestartPolicy                 string             `json:"restartPolicy"`
	TerminationGracePeriodSeconds int32              `json:"terminationGracePeriodSeconds"`
	DNSPolicy                     string             `json:"dnsPolicy"`
	ImagePullSecrets              []ImagePullSecrets `json:"imagePullSecrets"`
	SchedulerName                 string             `json:"schedulerName"`
}
