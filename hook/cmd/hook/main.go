package main

import (
	"github.com/emicklei/go-restful"
	"hook/internal/config"
	"k8s.io/klog/v2"
	"net/http"
)

func main() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	//Register
	config.Register(wsContainer)
	klog.Info("start listening on localhost:" + "8080")
	server := &http.Server{Addr: ":" + "8080", Handler: wsContainer}
	klog.Fatal(server.ListenAndServe())
}
