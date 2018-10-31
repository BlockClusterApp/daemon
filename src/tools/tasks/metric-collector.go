package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"log"
	"net/http"
)

type NodeMetricRequestEntity struct {
	NodeName string     `json:"nodeName"`
	Usage    dtos.Usage `json:"usage"`
}

type Container struct {
	Name  string     `json:"name"`
	Usage dtos.Usage `json:"usage"`
}

type PodMetricRequestEntity struct {
	PodName    string      `json:"podName"`
	Namespace  string      `json:"namespace"`
	Containers []Container `json:"containers"`
}

func getNodeMetrics() {
	path := "/apis/metrics.k8s.io/v1beta1/nodes"

	response, err := helpers.MakeKubeRequest(http.MethodGet, path, nil)

	if err != nil {
		helpers.GetLogger().Printf("Error fetching node metrics %s", err.Error())
		return
	}

	var nodeMetricResponse = dtos.NodeMetricResponse{}
	err = json.Unmarshal([]byte(response), nodeMetricResponse)

	if err != nil {
		helpers.GetLogger().Printf("Error unmarshalling node metrics %s", err.Error())
		return
	}

	if len(nodeMetricResponse.Items) == 0 {
		helpers.GetLogger().Printf("No node metric item")
		return
	}

	requestItems := make([]NodeMetricRequestEntity, len(nodeMetricResponse.Items))

	for i, item := range nodeMetricResponse.Items {
		requestItems[i].NodeName = item.Metadata.Name
		requestItems[i].Usage = item.Usage
	}

	jsonBody, err := json.Marshal(requestItems)

	if err != nil {
		helpers.GetLogger().Printf("Error converting request items to json %s", err.Error())
		return
	}

	requestObject := fmt.Sprintf("{\"nodes\": \"%s\"}", jsonBody)

	bc := helpers.GetBlockclusterInstance()
	bc.SendRequest("/metrics", requestObject)
}

func getPodMetrics() {
	paths := []string{
		"/apis/metrics.k8s.io/v1beta1/pods?labelSelector=app%3Dblockcluster-app",
		"/apis/metrics.k8s.io/v1beta1/pods?labelSelector=app%3Dhyperion",
		"/apis/metrics.k8s.io/v1beta1/namespaces/blockcluster/pods",
		"/apis/metrics.k8s.io/v1beta1/pods?labelSelector=appType%3Ddynamo",
	}

	for _, _path := range paths {
		go func(path string) {
			response, err := helpers.MakeKubeRequest(http.MethodGet, path, nil)

			if err != nil {
				helpers.GetLogger().Printf("Error fetching pod metrics %s | %s", path, err.Error())
				return
			}

			var podMetricResponse = dtos.PodMetricResponse{}
			err = json.Unmarshal([]byte(response), podMetricResponse)

			if err != nil {
				helpers.GetLogger().Printf("Error unmarshalling pod metrics %s", err.Error())
				return
			}

			if len(podMetricResponse.Items) == 0 {
				helpers.GetLogger().Printf("No pod metric item")
				return
			}

			requestItems := make([]PodMetricRequestEntity, len(podMetricResponse.Items))
			for i, item := range podMetricResponse.Items {
				requestItems[i].PodName = item.Metadata.Name
				requestItems[i].Namespace = item.Metadata.Namespace
				requestItems[i].Containers = make([]Container, len(item.Containers))
				for j, container := range item.Containers {
					requestItems[i].Containers[j].Name = container.Name
					requestItems[i].Containers[j].Usage = container.Usage
				}
			}

			jsonBody, err := json.Marshal(requestItems)

			if err != nil {
				helpers.GetLogger().Printf("Error converting request items to json %s", err.Error())
				return
			}

			requestObject := fmt.Sprintf("{\"pods\": \"%s\"}", jsonBody)
			bc := helpers.GetBlockclusterInstance()

			log.Printf("Pod metrics %s", path)
			bc.SendRequest("/metrics", requestObject)
		}(_path)
	}
}

func UpdateMetrics() {
	log.Printf("Metric collection started")
	go getNodeMetrics()
	go getPodMetrics()
}
