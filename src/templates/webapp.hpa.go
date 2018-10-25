package templates

func GetWebappHPATemplate() string {
  return `{
"apiVersion": "autoscaling/v1",
  "kind": "HorizontalPodAutoscaler",
  "metadata": {
    "name": "blockcluster-hpa",
    "namespace": "%__NAMESPACE__%"
  },
  "spec": {
    "maxReplicas": 10,
    "minReplicas": 2,
    "scaleTargetRef": {
      "apiVersion": "extensions/v1beta1",
      "kind": "Deployment",
      "name": "blockcluster-webapp-deploy"
    },
    "targetCPUUtilizationPercentage": 60
  }
}`
}