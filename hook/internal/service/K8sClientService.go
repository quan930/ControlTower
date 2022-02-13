package service

import v1 "github.com/quan930/ControlTower/ControlTower-operator/api/v1"

type K8sClientService interface {
	ListHook() (*v1.HookList, error)
}
