package k8s

import (
	scheduler "github.com/tioxy/scheduler/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k KubernetesAPI) CreateCronJob(sj scheduler.SimpleJob) error {
	cronJob := sj.ToCronJob()
	_, err := k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Create(cronJob)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) DeleteCronJob(sj scheduler.SimpleJob) error {
	cronJob := sj.ToCronJob()
	err := k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Delete(cronJob.ObjectMeta.Name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) UpdateCronJob(sj scheduler.SimpleJob) error {
	cronJob := sj.ToCronJob()
	_, err := k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Update(cronJob)

	if err != nil {
		return err
	}

	return nil
}
