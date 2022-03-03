package service

import (
	cloudv1 "github.com/quan930/ControlTower/ControlTower-operator/api/v1"
	v13 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HookService interface {
	GetJobByCheckGitEvent(event cloudv1.GitEvent, hook *cloudv1.Hook) (*v13.Job, *string)
	UpdateDeployment(imageEvent cloudv1.ImageEvent, hook *cloudv1.Hook, client client.Client)
}
