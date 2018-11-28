package templates

func GetWebappDeploymentTemplate() string {
	return `{
  "apiVersion": "apps/v1beta1",
  "kind": "Deployment",
  "metadata": {
    "labels": {
      "app": "blockcluster-app",
      "name": "blockcluster",
      "namespace": "%__NAMESPACE__%"
    },
    "name": "blockcluster-webapp-deploy"
  },
  "spec": {
    "replicas": 3,
    "selector": {
      "matchLabels": {
        "name": "blockcluster"
      }
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "blockcluster-app",
          "env": "production",
          "name": "blockcluster"
        }
      },
      "spec": {
        "containers": [
          {
            "env": [
              {
                "name": "WEB_ENV",
                "value": "production"
              },
              {
                "name": "NODE_ENV",
                "value": "production"
              },
              {
                "name": "NODE_ENV",
                "value": "production"
              },
              {
                "name": "MONGO_URL",
                "value": "%__MONGO_URL__%"
              },
              {
				"name": "ROOT_URL",
				"value": "%__ROOT_URL__%"
              },
              {
                "name": "REDIS_HOST",
                "value": "%__REDIS_HOST__%"
              },
              {
                "name": "REDIS_PORT",
                "value": "%__REDIS_PORT__%"
              },
              {
                "name": "NAMESPACE",
                "value": "%__NAMESPACE__%"
              },
              {
                "name": "RAZORPAY_ID",
				"value": "%__RAZORPAY_ID__%"
              },
              {
                "name": "RAZORPAY_KEY",
				"value": "%__RAZORPAY_KEY__%"
              },
              {
                "name": "KUBERNETES_NODE_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "spec.nodeName"
                  }
                }
              },
              {
                "name": "KUBERNETES_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              }
            ],
            "image": "%__IMAGE_URL__%",
            "imagePullPolicy": "IfNotPresent",
            "livenessProbe": {
              "exec": {
                "command": [
                  "cat",
                  "/tmp/webapp.lock"
                ]
              },
              "failureThreshold": 3,
              "initialDelaySeconds": 30,
              "periodSeconds": 15,
              "successThreshold": 1,
              "timeoutSeconds": 1
            },
            "name": "blockcluster",
            "ports": [
              {
                "containerPort": 3000,
                "name": "http",
                "protocol": "TCP"
              }
            ],
            "readinessProbe": {
              "failureThreshold": 3,
              "httpGet": {
                "path": "/ping",
                "port": 3000,
                "scheme": "HTTP"
              },
              "periodSeconds": 5,
              "successThreshold": 1,
              "timeoutSeconds": 1
            },
            "resources": {
              "limits": {
                "cpu": "600m",
                "memory": "1Gi"
              },
              "requests": {
                "cpu": "150m",
                "memory": "512Mi"
              }
            },
            "volumeMounts": [
              {
                "mountPath": "/tmp/logs/",
                "name": "webapp-logs"
              }
            ]
          }
        ],
        "dnsPolicy": "ClusterFirst",
        "imagePullSecrets": [
          {
            "name": "blockcluster-regsecret"
          }
        ],
        "restartPolicy": "Always",
        "volumes": [
          {
            "hostPath": {
              "path": "/tmp/webapp-logs/",
              "type": ""
            },
            "name": "webapp-logs"
          }
        ]
      }
    }
  }
}`
}
