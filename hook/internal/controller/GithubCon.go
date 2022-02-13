package controller

import (
	"github.com/emicklei/go-restful"
	"hook/internal/pojo"
	"hook/internal/service"
	"k8s.io/klog/v2"
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
		response.WriteEntity(pojo.NewResponse(500, "业务异常", nil).Body)
	} else {
		//klog.Info(pushPayload)
		klog.Info(pushPayload.Repository.URL)
		klog.Info(pushPayload.Ref)
		hooklist, err := k8sClientService.ListHook()
		klog.Info(hooklist)
		klog.Info(err)
		response.WriteEntity(pojo.NewResponse(200, "successful", nil).Body)
	}
}
