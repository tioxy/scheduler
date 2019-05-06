package k8s

import (
	scheduler "github.com/tioxy/scheduler/pkg"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k KubernetesAPI) CreateCronJob(sj scheduler.SimpleJob) error {
	cronJob := scheduler.ConvertSimpleJobToCronJob(sj)
	_, err := k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Create(&cronJob)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) DeleteCronJob(name string, namespace string) error {
	err := k.Client.BatchV1beta1().CronJobs(namespace).Delete(name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) FetchCronJob(name string, namespace string) (*batchv1beta1.CronJob, error) {
	cronJob, err := k.Client.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return cronJob, nil
}

func (k KubernetesAPI) UpdateCronJob(sj scheduler.SimpleJob) error {
	cronJob := scheduler.ConvertSimpleJobToCronJob(sj)
	_, err := k.Client.BatchV1beta1().CronJobs(cronJob.ObjectMeta.Namespace).Update(&cronJob)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) GetCronJobs(namespace string) ([]batchv1beta1.CronJob, error) {
	cronJobList, err := k.Client.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return cronJobList.Items, nil
}
