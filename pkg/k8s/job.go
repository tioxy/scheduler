package k8s

import (
	scheduler "github.com/tioxy/scheduler/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k KubernetesAPI) CreateJob(sj scheduler.SimpleJob) error {
	job := sj.ToJob()
	_, err := k.Client.BatchV1().Jobs(job.ObjectMeta.Namespace).Create(job)

	if err != nil {
		return err
	}

	return nil
}

func (k KubernetesAPI) DeleteJob(sj scheduler.SimpleJob) error {
	job := sj.ToJob()
	err = k.Client.BatchV1().Jobs(job.ObjectMeta.Namespace).Delete(job.ObjectMeta.Name, &metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}
