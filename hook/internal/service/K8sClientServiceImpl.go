package service

import (
	"context"
	"encoding/json"
	v1 "github.com/quan930/ControlTower/ControlTower-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type K8sClientServiceImpl struct {
}

func NewK8sClientServiceImpl() *K8sClientServiceImpl {
	initK8sClient()
	return &K8sClientServiceImpl{}
}

func (i K8sClientServiceImpl) ListHook() (*v1.HookList, error) {
	gvr := getGVR("cloud.lilqcn", "v1", "hooks")

	list, err := client.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var hookList v1.HookList
	if err := json.Unmarshal(data, &hookList); err != nil {
		return nil, err
	}

	return &hookList, nil
}

func (i K8sClientServiceImpl) UpdateHook(hook *v1.Hook) error {
	b, _ := json.Marshal(&hook)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	obj := &unstructured.Unstructured{Object: m}
	_, err := client.Resource(getGVR("cloud.lilqcn", "v1", "hooks")).Namespace(hook.GetNamespace()).UpdateStatus(ctx, obj, metav1.UpdateOptions{})
	return err
}

// getGVR :- gets GroupVersionResource for dynamic client
func getGVR(group, version, resource string) schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
}

var client dynamic.Interface
var ctx context.Context

func initK8sClient() {
	ctx = context.Background()
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	client, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
