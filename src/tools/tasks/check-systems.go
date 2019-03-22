package tasks

import "github.com/BlockClusterApp/daemon/src/helpers"

func CheckAllSystems() {
	helpers.GetLogger().Printf("GoCron: Checking systems")
	namespaces := helpers.GetNamespaces()

	for _, namespace := range namespaces {
		go helpers.CheckAndDeployHyperion(namespace)
	}
}