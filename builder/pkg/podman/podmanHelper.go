package podmanHelper

import (
	"context"
	"github.com/containers/buildah/define"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

//BuildImage build image
//dockerfileName eg:Dockerfile or Dockerfile2
//contextDirectory is the default source location for COPY and ADD commands. eg:"./temp"
func BuildImage(ctx context.Context, dockerfileName string, contextDirectory string, imageName string) error {
	_,err := images.Build(ctx, []string{contextDirectory+"/"+dockerfileName},entities.BuildOptions{
		BuildOptions: define.BuildOptions{
			Output: imageName,
			ContextDirectory: contextDirectory,
			NoCache: true,
		},
	})
	return err
}

func GetImageIDAndSha256(ctx context.Context, ImageName string) (string,string,error) {
	imageInspectReport, err := images.GetImage(ctx,ImageName,&images.GetOptions{})
	if err != nil {
		return "", "", err
	}
	return imageInspectReport.ID, imageInspectReport.Digest.Hex(),nil
}

func PushImage(ctx context.Context, registryUser string, registryPassword string, imageName string) error{
	return images.Push(ctx,imageName,imageName,&images.PushOptions{
		Username: &registryUser,
		Password: &registryPassword,
	})
}

func RemoveImage(ctx context.Context, imageName string) error {
	_,errs := images.Remove(ctx, []string{imageName},&images.RemoveOptions{})
	return errs[0]
}

func PullImage(ctx context.Context, imageName string) error {
	_, err := images.Pull(ctx, imageName, &images.PullOptions{})
	return err
}