package main

import (
	"builder/internal/git"
	"k8s.io/klog/v2"
)

func main() {
	//https://github.com/lianglitest/testimage

	repo := git.Clone("https://github.com/lianglitest/testimage", "./temp")
	err := git.Checkout(repo, "test")
	if err != nil {
		klog.Error(err)
	}
}
