package pkg

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

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

func (k KubernetesAPI) CreateJob(sj SimpleJob) error {
	job := sj.ToJob()
	_, err := k.Client.BatchV1().Jobs(job.ObjectMeta.Namespace).Create(job)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) CreateCronJob(sj SimpleJob) error {
	cronJob := sj.ToCronJob()
	_, err := k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Create(cronJob)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) DeleteJob(sj SimpleJob) error {
	job := sj.ToJob()
	err = k.Client.BatchV1().Jobs(job.ObjectMeta.Namespace).Delete(job.ObjectMeta.Name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) DeleteCronJob(sj SimpleJob) error {
	cronJob := sj.ToCronJob()
	err = k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Delete(cronJob.ObjectMeta.Name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) UpdateJob(sj SimpleJob) error {
	job := sj.ToJob()
	_, err = k.Client.BatchV1().Jobs(job.ObjectMeta.Namespace).Update(job)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) UpdateCronJob(sj SimpleJob) error {
	cronJob := sj.ToCronJob()
	_, err = k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Update(cronJob)

	if err != nil {
		return err
	}

	return nil
}

func areContainersEqual(srcContainers []v1.Container, destContainers []v1.Container) bool {
	if reflect.DeepEqual(srcContainers, destContainers) {
		return true
	}
	return false
}
