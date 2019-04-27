package pkg

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
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
	job := &batchv1.Job{
		ObjectMeta: sj.createObjectMeta(),
		Spec:       sj.createJobSpec(),
	}
	_, err := k.Client.BatchV1().Jobs(sj.Namespace).Create(job)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) CreateCronJob(sj SimpleJob) error {
	cronJob := &batchv1beta1.CronJob{
		ObjectMeta: sj.createObjectMeta(),
		Spec: batchv1beta1.CronJobSpec{
			Schedule: sj.Cron,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: sj.createJobSpec(),
			},
		},
	}
	_, err := k.Client.BatchV1beta1().CronJobs(sj.Namespace).Create(cronJob)

	if err != nil {
		return err
	}

	return nil
}
