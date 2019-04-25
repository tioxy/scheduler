package pkg

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var err error

type KubernetesAPI struct {
	Client kubernetes.Interface
}

func (k KubernetesAPI) CreateJob(sj SimpleJob) error {
	job := &batchv1.Job{
		ObjectMeta: sj.CreateObjectMeta(),
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers:    sj.Containers,
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &sj.MaxRetries,
		},
	}
	_, err := k.Client.BatchV1().Jobs(sj.Namespace).Create(job)

	if err != nil {
		return err
	}

	return nil
}
