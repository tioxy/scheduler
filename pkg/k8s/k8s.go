package k8s

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const allNamespaces = ""

var err error

type KubernetesAPI struct {
	Client kubernetes.Interface
}

func CreateKubernetesAPI() KubernetesAPI {
	config, err := rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	return KubernetesAPI{
		Client: clientSet,
	}
}

func areContainersEqual(srcContainers []v1.Container, destContainers []v1.Container) bool {
	if reflect.DeepEqual(srcContainers, destContainers) {
		return true
	}
	return false
}
