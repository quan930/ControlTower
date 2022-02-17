# ControlTower - builder
### version
+ 0.0.4
### test
```shell
docker run --privileged -e DOCKER_TLS_CERTDIR="" -d --name dockerd  docker:20.10.12-dind
docker run --rm -it --link dockerd:docker docker:20.10.12-git sh
docker run -d --link dockerd:docker -e REPO="https://github.com/lianglitest/testimage" -e BRANCH="main" -e DOCKERFILE="Dockerfile" -e IMAGE="lilqcn/testimage:1.7" -e USER="xxx" -e PASSWORD="xxx" lilqcn/builder:0.0.4
```