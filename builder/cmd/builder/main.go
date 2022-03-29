package main

import (
	"context"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/quan930/ControlTower/builder/pkg/git"
	podmanHelper "github.com/quan930/ControlTower/builder/pkg/podman"
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

	repo, err := git.Clone(repoURL, "./temp")
	if err != nil {
		klog.Fatal(err)
	}
	err = git.Checkout(repo, branch)
	if err != nil {
		klog.Fatal(err)
	}

	var conn context.Context
	conn, err = bindings.NewConnection(context.Background(), "unix:///usr/lib/systemd/system/podman.socket")
	if err != nil {
		klog.Warning(err)
		return
	}
	err = podmanHelper.BuildImage(conn,dockerfilePath,"./temp",image)
	if err != nil {
		klog.Fatal(err)
	}
	err = podmanHelper.PushImage(conn, username, password, image)
	if err != nil {
		klog.Fatal(err)
	}
}
