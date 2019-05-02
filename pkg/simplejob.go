package pkg

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
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
			"job":        sj.Name,
			"created-by": "scheduler",
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

func (sj SimpleJob) IsScheduled() bool {
	if sj.Cron == "" {
		return false
	}
	return true
}

func (sj SimpleJob) ToJob() *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: sj.createObjectMeta(),
		Spec:       sj.createJobSpec(),
	}
}

func (sj SimpleJob) ToCronJob() *batchv1beta1.CronJob {
	return &batchv1beta1.CronJob{
		ObjectMeta: sj.createObjectMeta(),
		Spec: batchv1beta1.CronJobSpec{
			Schedule: sj.Cron,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: sj.createJobSpec(),
			},
		},
	}
}
