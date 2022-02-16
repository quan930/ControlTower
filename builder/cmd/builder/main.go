package main

import (
	"builder/internal/git"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"k8s.io/klog/v2"
)

func main2() {
	//https://github.com/lianglitest/testimage

	repo := git.Clone("https://github.com/lianglitest/testimage", "./temp")
	err := git.Checkout(repo, "test")
	if err != nil {
		klog.Error(err)
	}
}

func main() {
	repo := git.Clone("https://github.com/lianglitest/testimage", "./temp")
	err := git.Checkout(repo, "test")
	if err != nil {
		klog.Error(err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	for true {

	}
}
