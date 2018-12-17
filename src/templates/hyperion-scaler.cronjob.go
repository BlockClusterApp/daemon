package templates

func GetHyperionScalerCronJobTemplate() string {
	return `{
  "apiVersion": "batch/v1beta1",
  "kind": "CronJob",
  "metadata": {
    "name": "hyperion-scaler"
  },
  "spec": {
    "concurrencyPolicy": "Forbid",
    "successfulJobsHistoryLimit": 0,
    "failedJobsHistoryLimit": 1,
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
                "imagePullPolicy": "Always",
                "env": [
                  {
                    "name": "K8_URI",
                    "value": "%__K8S_HOST__%"
                  },
                  {
                    "name": "K8_USER",
                    "value": "%__K8S_USER__%"
                  },
                  {
                    "name": "K8_PASS",
                    "value": "%__K8S_PASS__%"
                  },
                  {
                    "name": "K8_VERSION",
                    "value": "1.9"
                  },
                  {
                    "name": "PROM_NAMESPACE",
                    "value": "blockcluster-monitoring"
                  },
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
