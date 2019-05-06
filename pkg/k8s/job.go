package k8s

import (
	scheduler "github.com/tioxy/scheduler/pkg"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k KubernetesAPI) CreateJob(sj scheduler.SimpleJob) error {
	job := scheduler.ConvertSimpleJobToJob(sj)
	_, err := k.Client.BatchV1().Jobs(job.ObjectMeta.Namespace).Create(&job)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) DeleteJob(name string, namespace string) error {
	err = k.Client.BatchV1().Jobs(namespace).Delete(name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) FetchJob(name string, namespace string) (*batchv1.Job, error) {
	job, err := k.Client.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (k KubernetesAPI) ListJobs(namespace string) ([]batchv1.Job, error) {
	jobList, err := k.Client.BatchV1().Jobs(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return jobList.Items, nil
}
