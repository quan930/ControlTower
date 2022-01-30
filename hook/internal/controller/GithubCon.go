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

func (c GithubCon) Test(request *restful.Request, response *restful.Response)  {
	len := request.Request.ContentLength   // 获取请求实体长度
	body := make([]byte, len)  // 创建存放请求实体的字节切片
	request.Request.Body.Read(body)        // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
	str := string(body[:])
	klog.Info(str)
	response.WriteEntity(pojo.NewResponse(200, "successful", str).Body)
}