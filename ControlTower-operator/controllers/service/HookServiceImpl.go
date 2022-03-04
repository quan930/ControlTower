package service

import (
	"context"
	"github.com/google/uuid"
	cloudv1 "github.com/quan930/ControlTower/ControlTower-operator/api/v1"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/batch/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type HookServiceImpl struct {
}

func NewHookServiceImpl() *HookServiceImpl {
	return &HookServiceImpl{}
}

func (i HookServiceImpl) GetJobByCheckGitEvent(event cloudv1.GitEvent, hook *cloudv1.Hook) (*v13.Job, *string) {
	for _, item := range hook.Spec.Hooks {
		if event.GitRepository == item.GitRepository {
			for _, branch := range item.Branches {
				if branch == event.Branch {
					imageName := item.ImageRepository + ":" + time.Now().Format("20060102-1504")
					size1 := int32(1)
					size5 := int32(5)
					tr := true
					job := &v13.Job{
						ObjectMeta: metav1.ObjectMeta{
							Name:      hook.Name + "-buildimagejob" + "-" + uuid.New().String()[0:8],
							Namespace: "controltower-operator-system",
						},
						Spec: v13.JobSpec{
							Completions:  &size1,
							Parallelism:  &size1,
							BackoffLimit: &size5,
							Template: v12.PodTemplateSpec{
								Spec: v12.PodSpec{
									RestartPolicy: v12.RestartPolicy("OnFailure"),
									Volumes: []v12.Volume{{
										Name:         "lifecycle",
										VolumeSource: v12.VolumeSource{EmptyDir: &v12.EmptyDirVolumeSource{}},
									}},
									Containers: []v12.Container{{
										Image:           "lilqcn/builder:0.4.16-dind",
										Name:            "dind",
										Env:             []v12.EnvVar{{Name: "DOCKER_TLS_CERTDIR", Value: ""}},
										SecurityContext: &v12.SecurityContext{Privileged: &tr},
										VolumeMounts:    []v12.VolumeMount{{Name: "lifecycle", MountPath: "/lifecycle"}},
									}, {
										Image: "lilqcn/builder:0.4.16",
										Name:  "builder",
										Env: []v12.EnvVar{
											{Name: "DOCKER_HOST", Value: "tcp://localhost:2375"},
											{Name: "REPO", Value: item.GitRepository},
											{Name: "BRANCH", Value: event.Branch},
											{Name: "DOCKERFILE", Value: item.Dockerfile},
											{Name: "IMAGE", Value: imageName},
											{Name: "USER", Value: item.ImageRepoUser},
											{Name: "PASSWORD", Value: item.ImageRepoPassword},
										},
										VolumeMounts: []v12.VolumeMount{{Name: "lifecycle", MountPath: "/lifecycle"}},
									}},
								},
							},
						},
					}
					// Set Hook instance as the owner and controller
					//ctrl.SetControllerReference(hook, job, r.Scheme)
					return job, &imageName
				}
			}
		}
	}
	return nil, nil
}

func (i HookServiceImpl) UpdateDeployment(imageEvent cloudv1.ImageEvent, hook *cloudv1.Hook, client client.Client) {
	var workloads []cloudv1.Workload
	for _, item := range hook.Spec.Hooks {
		if item.ImageRepository == imageEvent.ImageRepository {
			workloads = append(workloads, item.Workloads...)
		}
	}
	if len(workloads) > 0 {
		ctx := context.Background()
		for _, workload := range workloads {
			newImage := imageEvent.ImageRepository + ":" + imageEvent.ImageTag
			if workload.Type == "Deployment" {
				updateDeployment(client, ctx, workload, newImage)
			} else if workload.Type == "StatefulSet" {
				//StatefulSet
				updateStatefulSet(client, ctx, workload, newImage)
			} else {
				//DaemonSet
				updateDaemonSet(client, ctx, workload, newImage)
			}
		}
	}
}

func updateDeployment(client client.Client, ctx context.Context, workload cloudv1.Workload, newImage string) {
	foundDeployment := &v1.Deployment{}
	err := client.Get(ctx, types.NamespacedName{Name: workload.Name, Namespace: workload.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("not found Deployment")
	}
	containers := &foundDeployment.Spec.Template.Spec.Containers
	for j, container := range *containers {
		if container.Name == workload.ContainerName {
			foundDeployment.Spec.Template.Spec.Containers[j].Image = newImage
		}
	}
	err = client.Update(ctx, foundDeployment)
	if err != nil {
		klog.Error(err, "Failed to update Deployment image")
	}
}

func updateDaemonSet(client client.Client, ctx context.Context, workload cloudv1.Workload, newImage string) {
	foundDaemonSet := &v1.DaemonSet{}
	err := client.Get(ctx, types.NamespacedName{Name: workload.Name, Namespace: workload.Namespace}, foundDaemonSet)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("not found DaemonSet")
	}
	containers := &foundDaemonSet.Spec.Template.Spec.Containers
	for j, container := range *containers {
		if container.Name == workload.ContainerName {
			foundDaemonSet.Spec.Template.Spec.Containers[j].Image = newImage
		}
	}
	err = client.Update(ctx, foundDaemonSet)
	if err != nil {
		klog.Error(err, "Failed to update DaemonSet image")
	}
}

func updateStatefulSet(client client.Client, ctx context.Context, workload cloudv1.Workload, newImage string) {
	foundStatefulSet := &v1.StatefulSet{}
	err := client.Get(ctx, types.NamespacedName{Name: workload.Name, Namespace: workload.Namespace}, foundStatefulSet)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("not found StatefulSet")
	}
	containers := &foundStatefulSet.Spec.Template.Spec.Containers
	for j, container := range *containers {
		if container.Name == workload.ContainerName {
			foundStatefulSet.Spec.Template.Spec.Containers[j].Image = newImage
		}
	}
	err = client.Update(ctx, foundStatefulSet)
	if err != nil {
		klog.Error(err, "Failed to update StatefulSet image")
	}
}
