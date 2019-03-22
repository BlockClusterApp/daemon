package templates

import "fmt"

func GetHyperionStatefulSet(diskSpace int32, apiVersion string) string {
	return fmt.Sprintf(`
{
  "apiVersion": "%s",
  "kind": "StatefulSet",
  "metadata": {
    "name": "hyperion",
    "labels": {
      "app": "hyperion"
    }
  },
  "spec": {
    "podManagementPolicy": "Parallel",
    "replicas": 3,
    "serviceName": "hyperion",
    "template": {
      "metadata": {
        "labels": {
          "app": "hyperion"
        }
      },
      "spec": {
        "imagePullSecrets": [
          {
            "name": "blockcluster-regsecret"
          }
        ],
        "containers": [
          {
            "name": "ipfs",
            "image": "ipfs/go-ipfs:v0.4.14",
            "ports": [
              {
                "containerPort": 5001,
                "name": "api"
              },
              {
                "containerPort": 8080,
                "name": "gateway"
              }
            ],
            "volumeMounts": [
              {
                "name": "ipfs-storage",
                "mountPath": "/data/ipfs"
              }
            ]
          },
          {
            "name": "ipfs-cluster",
            "image": "402432300121.dkr.ecr.us-west-2.amazonaws.com/ipfs-cluster:latest",
            "env": [
              {
                "name": "CLUSTER_SECRET",
                "value": "6161616161616161616161616161616161616161616161616161616161616161"
              },
              {
                "name": "IPFS_API",
                "value": "/ip4/127.0.0.1/tcp/5001"
              }
            ],
            "ports": [
              {
                "containerPort": 9094,
                "name": "api"
              },
              {
                "containerPort": 9095,
                "name": "gateway"
              },
              {
                "containerPort": 9096,
                "name": "tcp"
              }
            ],
            "volumeMounts": [
              {
                "name": "ipfs-storage",
                "mountPath": "/data/ipfs-cluster"
              }
            ]
          },
          {
            "name": "ipfs-cluster-sidecar",
            "image": "joshorig/ipfs-cluster-k8s-sidecar",
            "env": [
              {
                "name": "IPFSCLUSTER_SIDECAR_POD_LABELS",
                "value": "app=hyperion"
              },
              {
                "name": "POD_IP",
                "valueFrom": {
                  "fieldRef": {
                    "fieldPath": "status.podIP"
                  }
                }
              }
            ]
          }
        ]
      }
    },
    "volumeClaimTemplates": [
      {
        "metadata": {
          "name": "ipfs-storage"
        },
        "spec": {
          "accessModes": [
            "ReadWriteOnce"
          ],
          "resources": {
            "requests": {
              "storage": "%dGi"
            }
          }
        }
      }
    ]
  }
}
`, apiVersion, diskSpace)
}

func GetHyperionService(apiVersion string) string {
	return fmt.Sprintf(`
{
  "apiVersion": "%s",
  "kind": "Service",
  "metadata": {
    "name": "hyperion"
  },
  "spec": {
    "type": "NodePort",
    "ports": [
      {
        "port": 5001,
        "targetPort": 5001,
        "protocol": "TCP",
        "name": "api"
      },
      {
        "port": 8080,
        "targetPort": 8080,
        "name": "gateway"
      },
      {
        "port": 9094,
        "targetPort": 9094,
        "name": "cluster-api"
      },
      {
        "port": 9095,
        "targetPort": 9095,
        "name": "cluster-gateway"
      },
      {
        "port": 9096,
        "targetPort": 9096,
        "name": "cluster-tcp",
        "protocol": "TCP"
      }
    ],
    "selector": {
      "app": "hyperion"
    }
  }
}
`, apiVersion)
}

func GetHyperionClusterRole(apiVersion string) string {
	if apiVersion == "" {
		apiVersion = "rbac.authorization.k8s.io/v1"
	}
	return fmt.Sprintf(`
{
  "kind": "ClusterRole",
  "apiVersion": "%s",
  "metadata": {
    "name": "ipfs-rbac"
  },
  "rules": [
    {
      "apiGroups": [
        ""
      ],
      "resources": [
        "pods"
      ],
      "verbs": [
        "get",
        "list"
      ]
    }
  ]
}
`, apiVersion)
}

func GetHyperionClusterRoleBinding(apiVersion string, namespace string) string {
	if apiVersion == "" {
		apiVersion = "rbac.authorization.k8s.io/v1"
	}
	return fmt.Sprintf(`
{
  "kind": "ClusterRoleBinding",
  "apiVersion": "%s",
  "metadata": {
    "name": "ipfs-rbac-binding"
  },
  "subjects": [
    {
      "kind": "ServiceAccount",
      "name": "default",
      "namespace": "%s"
    }
  ],
  "roleRef": {
    "kind": "ClusterRole",
    "name": "ipfs-rbac",
    "apiGroup": "rbac.authorization.k8s.io"
  }
}
`, apiVersion, namespace)
}
