package controller

import (
	"github.com/emicklei/go-restful"
	"hook/internal/pojo"
	"k8s.io/klog/v2"
)

type GithubCon struct {
}

func NewGithubCon() *GithubCon {
	return &GithubCon{}
}

func (c GithubCon) Test(request *restful.Request, response *restful.Response) {
	pushPayload := new(pojo.PushPayload)
	err := request.ReadEntity(&pushPayload)
	if err != nil {
		klog.Warning(err)
		response.WriteEntity(pojo.NewResponse(500, "业务异常", nil).Body)
	} else {
		klog.Info(pushPayload)
		klog.Info(pushPayload.Repository.URL)
		klog.Info(pushPayload.Ref)
		response.WriteEntity(pojo.NewResponse(200, "successful", nil).Body)
	}
}
