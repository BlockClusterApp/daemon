package helpers

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"log"
)

type GroupAPIMapping struct {
	APIVersion string
	Path       string
}

var KubeAPIMapping = map[string]map[string]*GroupAPIMapping{
	"1.9": {
		"clusterroles": {
			APIVersion: "rbac.authorization.k8s.io/v1",
			Path:       "apis/rbac.authorization.k8s.io/v1",
		},
		"clusterrolebindings": {
			APIVersion: "rbac.authorization.k8s.io/v1",
			Path:       "apis/rbac.authorization.k8s.io/v1",
		},
		"deployments": {
			APIVersion: "apps/v1",
			Path:       "apis/apps/v1",
		},
		"statefulsets": {
			APIVersion: "apps/v1",
			Path:       "apis/apps/v1",
		},
		"services": {
			APIVersion: "v1",
			Path:       "api/v1",
		},
		"configmaps": {
			APIVersion: "v1",
			Path:       "api/v1",
		},
	},
}

func GetKubeAPIVersion(kubeVersion *dtos.KubeVersion, service string) GroupAPIMapping {
	var version string;
	if kubeVersion == nil {
		version = "1.9"
	} else {
		version = fmt.Sprintf("%s.%s", kubeVersion.Major, kubeVersion.Minor)
	}

	if version == "1.9" {
		group := KubeAPIMapping[version][service]
		if group == nil {
			log.Fatalf("Invalid service specified %s/%s", group, service)
		}
		return *group
	} else {
		return GroupAPIMapping{
			APIVersion: "api/v1",
			Path:       "api/v1",
		}
	}
}
