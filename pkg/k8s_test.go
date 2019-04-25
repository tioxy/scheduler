package pkg

import (
	"testing"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		MaxRetries: 4,
		Scheduled:  false,
		Containers: []v1.Container{
			v1.Container{
				Name: "pi",
				Image: "perl",
				Command: []string{"perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			}
		}
	}

	err := api.CreateJob(sj)

	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = k.Client.BatchV1().Jobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})

	if err != nil {
		t.Fatal(err.Error())
	}
}
