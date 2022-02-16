# ControlTower - builder
### version
+ 0.0.2
### test
```shell
docker run --privileged -e DOCKER_TLS_CERTDIR="" -d --name dockerd  docker:20.10.12-dind
docker run --rm -it --link dockerd:docker docker:20.10.12-git sh
docker run -d --link dockerd:docker lilqcn/builder:0.0.2
```