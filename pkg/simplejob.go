package pkg

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SimpleJob struct {
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
	Containers []v1.Container `json:"containers"`
	Scheduled  bool           `json:"scheduled"`
	MaxRetries int32          `json:"maxRetries"`
	Cron       string         `json:"cron,omitempty"`
}

func (sj SimpleJob) CreateObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      sj.Name,
		Namespace: sj.Namespace,
		Labels: map[string]string{
			"job": sj.Name,
		},
	}
}
