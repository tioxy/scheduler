package k8s

import (
	"testing"

	scheduler "github.com/tioxy/scheduler/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateCronJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := scheduler.SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		Cron:       "0 0 * * *",
		MaxRetries: 4,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}

	err := api.CreateCronJob(sj)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = api.Client.BatchV1beta1().CronJobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteCronJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := scheduler.SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		Cron:       "0 0 * * *",
		MaxRetries: 4,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}

	err = api.CreateCronJob(sj)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = api.DeleteCronJob(sj.Name, sj.Namespace)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = api.Client.BatchV1beta1().CronJobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})
	if err == nil {
		t.Fatal("K8S CronJob exists after deletion")
	}
}

func TestUpdateCronJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := scheduler.SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		Cron:       "0 0 * * *",
		MaxRetries: 4,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}

	err = api.CreateCronJob(sj)
	if err != nil {
		t.Fatal(err.Error())
	}

	newSj := scheduler.SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		Cron:       "0 12 * * *",
		MaxRetries: 1,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi-new",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(1000)"},
			},
		},
	}

	err = api.UpdateCronJob(newSj)
	if err != nil {
		t.Fatal(err.Error())
	}

	currentCronJob, err := api.Client.BatchV1beta1().CronJobs(sj.Namespace).Get(sj.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if !areContainersEqual(newSj.Containers, currentCronJob.Spec.JobTemplate.Spec.Template.Spec.Containers) {
		t.Fatal("Updated CronJob doesn't match new Containers from SimpleJob")
	}
}

func TestFetchCronJob(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	sj := scheduler.SimpleJob{
		Name:       "pi",
		Namespace:  "default",
		Cron:       "0 0 * * *",
		MaxRetries: 4,
		Containers: []v1.Container{
			v1.Container{
				Name:    "pi",
				Image:   "perl",
				Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
			},
		},
	}

	err := api.CreateCronJob(sj)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = api.FetchCronJob(
		sj.Name,
		sj.Namespace,
	)

	if err != nil {
		t.Fatal("Could not fetch CronJob because " + err.Error())
	}
}

func TestGetCronJobs(t *testing.T) {
	api := &KubernetesAPI{
		Client: fake.NewSimpleClientset(),
	}

	simpleJobs := []scheduler.SimpleJob{
		{
			Name:       "pi-1000",
			Namespace:  "default",
			Cron:       "0 0 * * *",
			MaxRetries: 4,
			Containers: []v1.Container{
				v1.Container{
					Name:    "pi",
					Image:   "perl",
					Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(1000)"},
				},
			},
		},
		{
			Name:       "pi-2000",
			Namespace:  "kube-system",
			Cron:       "0 0 * * *",
			MaxRetries: 4,
			Containers: []v1.Container{
				v1.Container{
					Name:    "pi",
					Image:   "perl",
					Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
				},
			},
		},
	}

	for _, sj := range simpleJobs {
		err := api.CreateCronJob(sj)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	cronJobsAllNamespaces, err := api.GetCronJobs(allNamespaces)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(cronJobsAllNamespaces) != len(simpleJobs) {
		t.Fatal("The amount of CronJobs created do not match the amount of CronJobs received from all namespaces")
	}

	cronJobsDefaultNamespace, err := api.GetCronJobs("default")
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(cronJobsDefaultNamespace) != 1 {
		t.Fatal("The amount of CronJobs created do not match the amount of CronJobs received from 'default' namespace")
	}
}
