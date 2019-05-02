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
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}

	err := api.CreateJob(sj)

	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = api.Client.BatchV1().Jobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		MaxRetries: 4,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}

	err = api.CreateJob(sj)

	if err != nil {
		t.Fatal(err.Error())
	}

	err = api.DeleteJob(sj)

	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = api.Client.BatchV1().Jobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})

	if err == nil {
		t.Fatal("K8S Job exists after deletion")
	}
}

func TestUpdateJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		MaxRetries: 4,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}
	err = api.CreateJob(sj)
	if err != nil {
		t.Fatal(err.Error())
	}

	newSj := SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		MaxRetries: 1,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi-new",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(1000)"},
			},
		},
	}
	err = api.UpdateJob(newSj)
	if err != nil {
		t.Fatal(err.Error())
	}

	currentJob, err := api.Client.BatchV1().Jobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if !areContainersEqual(newSj.Containers, currentJob.Spec.Template.Spec.Containers) {
		t.Fatal("Updated Job doesn't match new Containers from SimpleJob")
	}
}
