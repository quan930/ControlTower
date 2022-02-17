package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"k8s.io/klog/v2"
	"os"
)

func Clone(gitUrl string, workPath string) *git.Repository {
	klog.Info("git clone ", gitUrl)

	repository, err := git.PlainClone(workPath, false, &git.CloneOptions{
		URL:      gitUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		klog.Fatal(err)
		return nil
	}
	return repository
}

func Checkout(repo *git.Repository, branchName string) error {
	klog.Info("git fetch")
	err := repo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info("git checkout ", branchName)
	wt, err := repo.Worktree()
	if err != nil {
		klog.Fatal(err)
	}
	return wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName)),
		Force:  true,
	})
	return nil
}
