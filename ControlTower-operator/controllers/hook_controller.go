/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/google/uuid"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/batch/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"time"

	cloudv1 "github.com/quan930/ControlTower/ControlTower-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HookReconciler reconciles a Hook object
type HookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Hook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
//+kubebuilder:rbac:groups=cloud.lilqcn,resources=hooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.lilqcn,resources=hooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.lilqcn,resources=hooks/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=jobs,verbs=get;list;watch;create;update;patch;delete
func (r *HookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("========== start ===============>")
	// 获取 MyBook 实例
	hook := &cloudv1.Hook{}
	ctx = context.Background()

	err := r.Get(ctx, req.NamespacedName, hook)
	if err != nil {
		if errors.IsNotFound(err) {
			// 对象未找到
			klog.Info("hook resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		klog.Error(err, "Failed to get Hook")
		return ctrl.Result{}, err
	}
	klog.Info("hook:", hook)
	//判断是否有git event
	if len(hook.Status.GitEvents) > 0 {
		klog.Info("GitEvents > 0")
		for i, event := range hook.Status.GitEvents {
			job, imagename := r.checkGitEvent(event, hook)
			if job != nil {
				klog.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
				err = r.Create(ctx, job)
				if err != nil {
					klog.Error(err, "Failed to create new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
					return ctrl.Result{}, err
				}
				hook.Status.GitEventHistory = append(hook.Status.GitEventHistory, cloudv1.GitEventHistory{GitRepository: event.GitRepository, Branch: event.Branch, DateTime: time.Now().Format("2006-01-02-15:04:05"), Status: "Running", BuildImageJob: job.Name, ImageName: *imagename})
				hook.Status.GitEvents = append(hook.Status.GitEvents[:i], hook.Status.GitEvents[i+1:]...)
			} else {
				hook.Status.GitEvents = append(hook.Status.GitEvents[:i], hook.Status.GitEvents[i+1:]...)
				hook.Status.GitEventHistory = append(hook.Status.GitEventHistory, cloudv1.GitEventHistory{GitRepository: event.GitRepository, Branch: event.Branch, DateTime: time.Now().Format("2006-01-02-15:04:05"), Status: "no need push"})
			}
			err = r.Status().Update(ctx, hook)
			if err != nil {
				klog.Error(err, "Failed to update Hook")
				return ctrl.Result{}, err
			}
			klog.Info(hook)
			return ctrl.Result{}, nil
		}
	}
	//判断是否有image event
	//todo 判断是否有image event
	if len(hook.Status.ImageEvents) > 0 {
		klog.Info("ImageEvents > 0")
		for i, event := range hook.Status.ImageEvents {
			//todo get update deployment
			deploys := getNeedUpdateDeployment(event, hook)
			if len(*deploys) > 0 {
				for _, deploy := range *deploys {
					foundDeployment := &v1.Deployment{}
					err = r.Get(ctx, types.NamespacedName{Name: deploy.DeployName, Namespace: deploy.Namespace}, foundDeployment)
					if err != nil && errors.IsNotFound(err) {
						klog.Info("not found deployment")
					}
					containers := &foundDeployment.Spec.Template.Spec.Containers
					for j, container := range *containers {
						if container.Name == deploy.ContainerName {
							foundDeployment.Spec.Template.Spec.Containers[j].Image = event.ImageRepository + ":" + event.ImageTag
						}
					}
					err = r.Update(ctx, foundDeployment)
					if err != nil {
						klog.Error(err, "Failed to update Hook")
					}
				}
			}
			hook.Status.ImageEvents = append(hook.Status.ImageEvents[:i], hook.Status.ImageEvents[i+1:]...)
			hook.Status.ImageEventHistory = append(hook.Status.ImageEventHistory, cloudv1.ImageEventHistory{ImageRepository: event.ImageRepository, ImageTag: event.ImageTag, DateTime: time.Now().Format("2006-01-02-15:04:05")})
			image := event.ImageRepository + ":" + event.ImageTag
			for j, history := range hook.Status.GitEventHistory {
				if history.ImageName == image {
					hook.Status.GitEventHistory[j].Status = "Successful" // Completed
					err = r.Update(ctx, hook)
					if err != nil {
						klog.Error(err, "Failed to update Hook")
						return ctrl.Result{}, err
					}
					klog.Info(hook)
					return ctrl.Result{}, nil
				}
			}
			err = r.Status().Update(ctx, hook)
			if err != nil {
				klog.Error(err, "Failed to update Hook")
				return ctrl.Result{}, err
			}
			klog.Info(hook)
			return ctrl.Result{}, nil
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//控制器监视的资源
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1.Hook{}).
		//将 Deployments 类型指定为要监视的辅助资源。对于每个部署类型的添加/更新/删除事件，事件处理程序会将每个事件映射到Request部署所有者的协调
		Owns(&v1.Deployment{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

//checkGitEvent 返回job image name
func (r *HookReconciler) checkGitEvent(event cloudv1.GitEvent, hook *cloudv1.Hook) (*v13.Job, *string) {
	for _, item := range hook.Spec.Hooks {
		if event.GitRepository == item.GitRepository {
			for _, branch := range item.Branches {
				if branch == event.Branch {
					klog.Info("need to deploy buildImage job")
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
										Image:           "lilqcn/builder:0.0.4-dind",
										Name:            "dind",
										Env:             []v12.EnvVar{{Name: "DOCKER_TLS_CERTDIR", Value: ""}},
										SecurityContext: &v12.SecurityContext{Privileged: &tr},
										VolumeMounts:    []v12.VolumeMount{{Name: "lifecycle", MountPath: "/lifecycle"}},
									}, {
										Image: "lilqcn/builder:0.0.4",
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
					ctrl.SetControllerReference(hook, job, r.Scheme)
					return job, &imageName
				}
			}
		}
	}
	return nil, nil
}

//checkJob 检验image job
//func (r *HookReconciler) checkJob(hook *cloudv1.Hook) bool {
//	for _, event := range hook.Spec.ImageEvents {
//		image := event.ImageRepository + ":" + event.ImageTag
//		for i, history := range hook.Spec.GitEventHistory {
//			if history.ImageName == image {
//				hook.Spec.GitEventHistory[i].Status = "Successful" // Completed
//				return true
//			}
//		}
//	}
//	return false
//}

func getNeedUpdateDeployment(imageEvent cloudv1.ImageEvent, hook *cloudv1.Hook) *[]cloudv1.Deploy {
	var deploys []cloudv1.Deploy
	for _, item := range hook.Spec.Hooks {
		if item.ImageRepository == imageEvent.ImageRepository {
			deploys = append(deploys, item.Deploys...)
		}
	}
	return &deploys
}
