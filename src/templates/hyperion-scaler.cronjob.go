package templates

func GetHyperionScalerCronJobTemplate() string {
	return `{
  "apiVersion": "batch/v1beta1",
  "kind": "CronJob",
  "metadata": {
    "name": "hyperion-scaler"
  },
  "spec": {
    "schedule": "*/5 * * * *",
    "jobTemplate": {
      "spec": {
        "template": {
          "spec": {
            "restartPolicy": "OnFailure",
            "serviceAccountName": "hyperion-scaler-sa",
            "containers": [
              {
                "name": "scaler",
                "image": "402432300121.dkr.ecr.us-west-2.amazonaws.com/hyperion-scaler:latest",
                "imagePullPolicy": "IfNotPresent",
                "env": [
                  {
                    "name": "NAMESPACE",
                    "value": "%__NAMESPACE__%"
                  }
                ]
              }
            ],
            "imagePullSecrets": [
              {
                "name": "blockcluster-regsecret"
              }
            ]
          }
        }
      }
    }
  }
}`
}
