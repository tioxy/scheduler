package pkg

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SimpleJob struct {
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
	Containers []v1.Container `json:"containers"`
	MaxRetries int32          `json:"maxRetries"`
	Cron       string         `json:"cron,omitempty"`
}

func (sj SimpleJob) createObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      sj.Name,
		Namespace: sj.Namespace,
		Labels: map[string]string{
			"job": sj.Name,
		},
	}
}

func (sj SimpleJob) createJobSpec() batchv1.JobSpec {
	return batchv1.JobSpec{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Containers:    sj.Containers,
				RestartPolicy: v1.RestartPolicyNever,
			},
		},
		BackoffLimit: &sj.MaxRetries,
	}
}
