package pkg

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

func ConvertJobToSimpleJob(job batchv1.Job) SimpleJob {
	return SimpleJob{
		Name:       job.ObjectMeta.Name,
		Namespace:  job.ObjectMeta.Namespace,
		Containers: job.Spec.Template.Spec.Containers,
		MaxRetries: *job.Spec.BackoffLimit,
	}
}

func ConvertCronJobToSimpleJob(cronJob batchv1beta1.CronJob) SimpleJob {
	return SimpleJob{
		Name:       cronJob.ObjectMeta.Name,
		Namespace:  cronJob.ObjectMeta.Namespace,
		Containers: cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers,
		MaxRetries: *cronJob.Spec.JobTemplate.Spec.BackoffLimit,
		Cron:       cronJob.Spec.Schedule,
	}
}

func ConvertSimpleJobToCronJob(sj SimpleJob) batchv1beta1.CronJob {
	return batchv1beta1.CronJob{
		ObjectMeta: sj.generateObjectMeta(),
		Spec: batchv1beta1.CronJobSpec{
			Schedule: sj.Cron,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: sj.generateObjectMeta(),
				Spec:       sj.generateJobSpec(),
			},
		},
	}
}

func ConvertSimpleJobToJob(sj SimpleJob) batchv1.Job {
	return batchv1.Job{
		ObjectMeta: sj.generateObjectMeta(),
		Spec:       sj.generateJobSpec(),
	}
}
