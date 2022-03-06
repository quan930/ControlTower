package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	"io"
	"k8s.io/klog/v2"
	"os"
)

func BuildImage(cli *client.Client, dockerfilePath string, filePath string, imageName string) error {
	imageBuildResponse, err := cli.ImageBuild(context.Background(), getContext(filePath), types.ImageBuildOptions{
		Dockerfile: dockerfilePath,
		Tags:       []string{imageName},
	})
	if err != nil {
		klog.Error(err)
		return err
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		klog.Error(err, " :unable to read image build response")
		return err
	}
	return nil
}

func PushImage(cli *client.Client, registryUser string, registryPassword string, image string) error {
	authConfig := types.AuthConfig{
		Username: registryUser,
		Password: registryPassword,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		klog.Error(err)
		return err
	}
	klog.Infof("Push docker image registry: %v %v", registryUser, registryPassword)

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	out, err := cli.ImagePush(context.TODO(), image, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		klog.Error(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		klog.Error(err, " :unable to read image build response")
		return err
	}
	return nil
}

func GetImageSha256(cli *client.Client, image string) (error,*string) {
	dis,err := cli.DistributionInspect(context.Background(),image,"")
	if err != nil {
		klog.Error(err)
		return err,nil
	}
	sha256 := dis.Descriptor.Digest.String()
	return nil,&sha256
}

func getContext(filePath string) io.Reader {
	// Use homedir.Expand to resolve paths like '~/repos/myrepo'
	fileP, _ := homedir.Expand(filePath)
	ctx, _ := archive.TarWithOptions(fileP, &archive.TarOptions{})
	return ctx
}
