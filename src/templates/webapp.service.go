package templates

func GetWebappServiceTemplate() string {
  return `{
"apiVersion": "v1",
"kind": "Service",
  "metadata": {
    "labels": {
      "app": "blockcluster-app",
    },
    "name": "blockcluster-svc",
    "namespace": "%__NAMESPACE__%"
  },
  "spec": {
    "ports": [
      {
        "name": "http",
        "port": 80,
        "protocol": "TCP",
        "targetPort": 3000
      },
      {
        "name": "https",
        "port": 443,
        "protocol": "TCP",
        "targetPort": 3000
      }
    ],
    "selector": {
      "app": "blockcluster-app"
    },
    "sessionAffinity": "None",
    "type": "LoadBalancer"
  }
}`
}