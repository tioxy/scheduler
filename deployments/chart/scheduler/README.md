# Scheduler - Helm Chart

Standard chart definition to deploy Scheduler API on Kubernetes.

## [Configuration](#configuration)
-----

| **Parameter**                 | **Description**                                                                                                 	                | **Type**  	         |
|-----------------------------  |---------------------------------------------------------------------------------------------------------------------------------  |----------------------  |
| **`replicaCount`**           	| Amount of replicas from Kubernetes Deployment                                                                     	            | int       	         |
| **`image.repository`**        | Docker Image name (including repo if exists)                                                                     	                | str        	         |
| **`image.tag`**               | Docker Image tag                                                                         	                                        | str       	         |
| **`image.pullPolicy`**        | [How kubernetes pulls the image](https://kubernetes.io/docs/concepts/containers/images/#updating-images)                          | str       	         |
| **`image.port`**              | Container port and Environment variable **PORT** inside Kubernetes Deployment                                    	                | int       	         |
| **`image.debug`**             | If **true**, enables debug mode for application logs                                                             	                | bool       	         |
| **`service.type`**            | [Kubernetes Service type](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types)     | str        	         |
| **`service.port `**           | Exposed Port from the Kubernetes Service                                                                                     	    | int        	         |
| **`ingress.enabled`**         | If **true**, creates an ingress for the Kubernetes Service                                                                   	    | bool       	         |
| **`ingress.annotations`**     | Annotations for Kubernetes Ingress metadata                                                                          	            | map[str]str 	         |
| **`ingress.hosts`**           | Similar to a IngressRule from Kubernetes, but without `http` key                                                     	            | []IngressRule          |
| **`resources`**               | [Resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/) from Scheduler container       | ResourceRequirements   |
| **`rbac.enabled `**           | If **true**, create ServiceAccount, ClusterRole and ClusterRoleBinding for Scheduler to manage jobs                               | bool                   |


## [Example](#example)
-----

#### yourvalues.yaml
```yaml
# Example Helm Values file
replicaCount: 1

image:
  repository: tioxy/scheduler
  tag: latest
  pullPolicy: Always
  port: 8080
  debug: false

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: true
  annotations:
    your.special/annotation: value
  hosts:
    - host: scheduler.tioxy.com
      paths:
        - "/"
    - host: schedule.here.com
      paths:
        - "/"

resources:
    requests:
        memory: "64Mi"
        cpu: "100m"
    limits:
        memory: "256Mi"
        cpu: "500m"

rbac:
  enabled: true
```
