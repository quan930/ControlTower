package main

import (
	"github.com/docker/docker/client"
	"github.com/quan930/ControlTower/builder/pkg/docker"
	"github.com/quan930/ControlTower/builder/pkg/git"
	"k8s.io/klog/v2"
	"os"
)

func parseENV() (string, string, string, string, string, string) {
	var repo string
	var branch string
	var dockerfilePath string
	var image string
	var username string
	var password string
	repo = os.Getenv("REPO")
	branch = os.Getenv("BRANCH")
	dockerfilePath = os.Getenv("DOCKERFILE")
	image = os.Getenv("IMAGE")
	username = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	return repo, branch, dockerfilePath, image, username, password
}

func main() {
	repoURL, branch, dockerfilePath, image, username, password := parseENV()
	klog.Info("repo:", repoURL, "\nbranch:", branch, "\ndockerfilePath:", dockerfilePath, "\nimage", image, "\nuser", username, "\npassword", password)

	repo,err := git.Clone(repoURL, "./temp")
	if err != nil {
		klog.Fatal(err)
	}
	err = git.Checkout(repo, branch)
	if err != nil {
		klog.Fatal(err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		klog.Fatal(err)
	}

	err = docker.BuildImage(cli, dockerfilePath, "./temp", image)
	if err != nil {
		klog.Fatal(err)
	}
	err = docker.PushImage(cli, username, password, image)
	if err != nil {
		klog.Fatal(err)
	}

	f, err := os.Create("/lifecycle/main-terminated")
	defer f.Close()
	if err != nil {
		klog.Fatal(err)
	} else {
		klog.Info("/lifecycle/main-terminated")
	}
}
