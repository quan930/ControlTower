package controller

import (
	"github.com/emicklei/go-restful"
	v1 "github.com/quan930/ControlTower/ControlTower-operator/api/v1"
	"hook/internal/pojo"
	"hook/internal/service"
	"k8s.io/klog/v2"
	"strings"
)

type GithubCon struct {
}

func NewGithubCon() *GithubCon {
	k8sClientService = service.NewK8sClientServiceImpl()
	klog.Info("k8sClientService init")
	return &GithubCon{}
}

var k8sClientService service.K8sClientService

func (c GithubCon) GithubHook(request *restful.Request, response *restful.Response) {
	pushPayload := new(pojo.PushPayload)
	err := request.ReadEntity(&pushPayload)
	if err != nil {
		klog.Warning(err)
		response.WriteEntity(pojo.NewResponse(500, "read entity error", nil).Body)
	} else {
		//klog.Info(pushPayload)
		klog.Info(pushPayload.Repository.URL)
		klog.Info(pushPayload.Ref)
		hookList, err := k8sClientService.ListHook()
		if err != nil {
			klog.Info(err)
			response.WriteEntity(pojo.NewResponse(500, "update error", nil).Body)
			return
		}
		//klog.Info(hookList)
		hook := containGithub(pushPayload.Repository.URL, hookList)
		if hook == nil {
			response.WriteEntity(pojo.NewResponse(500, "webhook error", nil).Body)
		} else {
			hook.Status.GitEvents = append(hook.Status.GitEvents, v1.GitEvent{
				GitRepository: pushPayload.Repository.URL,
				Branch:        pushPayload.Ref[strings.LastIndex(pushPayload.Ref, "/")+1:],
			})
			klog.Info("name:" + hook.ObjectMeta.Name + "\tnamespace:" + hook.ObjectMeta.Namespace)
			err = k8sClientService.UpdateHook(hook)
			if err != nil {
				klog.Info(err)
				response.WriteEntity(pojo.NewResponse(500, "update error", nil).Body)
			} else {
				response.WriteEntity(pojo.NewResponse(200, "successful", nil).Body)
			}
		}
	}
}

func (c GithubCon) DockerhubHook(request *restful.Request, response *restful.Response) {
	buildPayload := new(pojo.BuildPayload)
	err := request.ReadEntity(&buildPayload)
	if err != nil {
		klog.Warning(err)
		response.WriteEntity(pojo.NewResponse(500, "read entity error", nil).Body)
	} else {
		//klog.Info(buildPayload)
		klog.Info(buildPayload.PushData.Tag)
		klog.Info(buildPayload.Repository.RepoName)
		hookList, err := k8sClientService.ListHook()
		if err != nil {
			klog.Info(err)
			response.WriteEntity(pojo.NewResponse(500, "update error", nil).Body)
			return
		}
		hook := containDockerHub(buildPayload.Repository.RepoName, hookList)
		if hook == nil {
			response.WriteEntity(pojo.NewResponse(500, "webhook error", nil).Body)
		} else {
			hook.Status.ImageEvents = append(hook.Status.ImageEvents, v1.ImageEvent{
				ImageRepository: buildPayload.Repository.RepoName,
				ImageTag:        buildPayload.PushData.Tag,
			})
			klog.Info("name:" + hook.ObjectMeta.Name + "\tnamespace:" + hook.ObjectMeta.Namespace)
			err = k8sClientService.UpdateHook(hook)
			if err != nil {
				klog.Info(err)
				response.WriteEntity(pojo.NewResponse(500, "update error", nil).Body)
			} else {
				response.WriteEntity(pojo.NewResponse(200, "successful", nil).Body)
			}
		}
	}
}

func containGithub(url string, hookList *v1.HookList) *v1.Hook {
	for _, item := range hookList.Items {
		for _, hook := range item.Spec.Hooks {
			if hook.GitRepository == url {
				return &item
			}
		}
	}
	return nil
}

func containDockerHub(repoName string, hookList *v1.HookList) *v1.Hook {
	for _, item := range hookList.Items {
		for _, hook := range item.Spec.Hooks {
			if hook.ImageRepository == repoName {
				return &item
			}
		}
	}
	return nil
}
