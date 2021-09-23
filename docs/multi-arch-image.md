# 多架构镜像构建指南

## 背景

为了适配越来越多的异构需求，程序的镜像构建不再满足于 linux/amd64 镜像，开始期望对于一个程序，能构建出适用于不同架构操作系统的镜像。

## 说明

最基本的多架构镜像，包含 linux/amd64 和 linux/arm64 镜像适配，至于 windows/amd64，linux/386，linux/mips64le 等镜像的适配，需要经过充分的测试。本文主要关注 linux/amd64 和 linux/arm64 的多架构镜像构建。

## 什么是多架构镜像？

以 golang:1.16 为例，它是一个典型的多架构镜像

```json
➜  ~ docker manifest inspect golang:1.16                                               
{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
   "manifests": [
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1795,
         "digest": "sha256:448a13037d13401ad9b31fabf91d6b8a3c5c35d336cc8af760b2ab4ed85d4155",
         "platform": {
            "architecture": "amd64",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1795,
         "digest": "sha256:6f662572d1418bd32d87ff8f6fe0e87aa4bbff74b6468f6d90d8df0202c9034c",
         "platform": {
            "architecture": "arm",
            "os": "linux",
            "variant": "v5"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1795,
         "digest": "sha256:e0bf5671a79a4f7766afa245f3ccdfc93a0509a1e01bfb9bfbf79742dbd171a4",
         "platform": {
            "architecture": "arm",
            "os": "linux",
            "variant": "v7"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1794,
         "digest": "sha256:c47161a791cacd4176f63b13feddf6923ad21f21026813f74aebe3c203800406",
         "platform": {
            "architecture": "arm64",
            "os": "linux",
            "variant": "v8"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1796,
         "digest": "sha256:633dff80ab7628865ee68d090ccc2fe52c7b6b9c9e34dcb7250f4efcfbc98dd0",
         "platform": {
            "architecture": "386",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1795,
         "digest": "sha256:92ef88c101f81bcae7a73cd76b3f6be76b9bc0d23291e2e4e288c159f5e94857",
         "platform": {
            "architecture": "mips64le",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1795,
         "digest": "sha256:f497dc3e6aaa05ffe3bd03b81ddd6f5e58c56b65ee36658531bc563aa03badbd",
         "platform": {
            "architecture": "ppc64le",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 1795,
         "digest": "sha256:08d14f280073634e7daecb8f3e0284227a38ddda2895b26261a92ed09f29b10c",
         "platform": {
            "architecture": "s390x",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 3401,
         "digest": "sha256:528d9a73e9f01c44bd725eb96c34d7f4505461ffc4201ed688ea5d7787ab8508",
         "platform": {
            "architecture": "amd64",
            "os": "windows",
            "os.version": "10.0.17763.2114"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 3323,
         "digest": "sha256:d2d5513c5c544adaffb3be2c0f21cec74f453186a8360feaa563770deb7b20ba",
         "platform": {
            "architecture": "amd64",
            "os": "windows",
            "os.version": "10.0.14393.4583"
         }
      }
   ]
}
```

### image manifest

可以看到 golang:1.16 由多个不同的 manifest 组成，每个 manifest 包含唯一的 digest，platform 中含有 os 和 arch 的信息，也就是说，本质上，golang:1.16 多架构镜像由多个不同 platfrom 的 manifest 组成，你可以使用 `docker pull golang@{digest}` 这种方式，拉取指定的 manifest

### target osarch

当你在不同架构的机器上执行`docker pull golang:1.16`时，docker 会做以下操作：

1. 获取当前机器的 os 和 arch 信息作为 target osarch
2. 使用 target osarch，去镜像仓库中拉取对应的 digest

这样，同一个 golang:1.16 镜像在不同架构的机器上执行 docker pull 的时候，都能拉取到对应的镜像

## 如何构建多架构镜像？

### 方式一：使用 docker manifests 组合多个已有的不同架构镜像为一个多平台镜像

开启 docker manifest 特性

```bash
→ vim ~/.docker/config.json
{
  "experimental": "enabled"
}

→ vim /etc/docker/daemon.json
{
  "experimental": true
}

→ systemctl daemon-reload
→ service docker restart
→ docker manifest --help
```

将 registry 中的不同架构镜像合为一个镜像

```bash
→ docker manifest create hub.c.163.com/kubecube/cube:v1.0.0-multi-arch hub.c.163.com/kubecube/cube:v1.0.0-arm64 hub.c.163.com/kubecube/cube:v1.0.0-amd64
```

将合成的多平台镜像推送到仓库

```bash
→ docker manifest push hub.c.163.com/kubecube/cube:v1.0.0-multi-arch 
```

### 方式二：使用 docker buildx 构建

`buildx`是 docker 的多平台镜像构建插件，其本质是翻译不同的指令集，并在此之上进行构建。要想使用 `buildx`，首先要确保 Docker 版本不低于 `19.03`，同时还要通过设置环境变量 `DOCKER_CLI_EXPERIMENTAL` 来启用。可以通过下面的命令来为当前终端启用 buildx 插件：

```bash
→ export DOCKER_CLI_EXPERIMENTAL=enabled
```

验证是否开启：

```bash
→ docker buildx version
github.com/docker/buildx v0.3.1-tp-docker 6db68d029599c6710a32aa7adcba8e5a344795a7
```

如果在某些系统上设置环境变量 DOCKER_CLI_EXPERIMENTAL 不生效（比如 Arch Linux）,你可以选择从源代码编译：

```bash
→ export DOCKER_BUILDKIT=1
→ docker build --platform=local -o . git://github.com/docker/buildx
→ mkdir -p ~/.docker/cli-plugins && mv buildx ~/.docker/cli-plugins/docker-buildx
```

启用 binfmt_misc

如果你使用的是 Docker 桌面版（MacOS 和 Windows），默认已经启用了 binfmt_misc，可以跳过这一步。

如果你使用的是 Linux，需要手动启用 binfmt_misc。大多数 Linux 发行版都很容易启用，不过还有一个更容易的办法，直接运行一个特权容器，容器里面写好了设置脚本：
```
→ docker run --rm --privileged docker/binfmt:66f9012c56a8316f9244ffd7622d7c21c1f6f28d
```
建议将 Linux 内核版本升级到 4.x 以上，特别是 CentOS 用户，你可能会遇到错误。

验证是 binfmt_misc 否开启：
```bash
→ ls -al /proc/sys/fs/binfmt_misc/
总用量 0
总用量 0
-rw-r--r-- 1 root root 0 11月 18 00:12 qemu-aarch64
-rw-r--r-- 1 root root 0 11月 18 00:12 qemu-arm
-rw-r--r-- 1 root root 0 11月 18 00:12 qemu-ppc64le
-rw-r--r-- 1 root root 0 11月 18 00:12 qemu-s390x
--w------- 1 root root 0 11月 18 00:09 register
-rw-r--r-- 1 root root 0 11月 18 00:12 status
```
验证是否启用了相应的处理器：
```bash
→ cat /proc/sys/fs/binfmt_misc/qemu-aarch64
enabled
interpreter /usr/bin/qemu-aarch64
flags: OCF
offset 0
magic 7f454c460201010000000000000000000200b7
mask ffffffffffffff00fffffffffffffffffeffff
```
从默认的构建器切换到多平台构建器
Docker 默认会使用不支持多 CPU 架构的构建器，我们需要手动切换。
先创建一个新的构建器：
```bash
→ docker buildx create --use --name mybuilder
```
启动构建器：
```bash 
→ docker buildx inspect mybuilder --bootstrap

[+] Building 5.0s (1/1) FINISHED
 => [internal] booting buildkit                                                                                                                          5.0s
 => => pulling image moby/buildkit:buildx-stable-1                                                                                                       4.4s
 => => creating container buildx_buildkit_mybuilder0                                                                                                     0.6s
Name:   mybuilder
Driver: docker-container

Nodes:
Name:      mybuilder0
Endpoint:  unix:///var/run/docker.sock
Status:    running
Platforms: linux/amd64, linux/arm64, linux/ppc64le, linux/s390x, linux/386, linux/arm/v7, linux/arm/v6
```
查看当前使用的构建器及构建器支持的 CPU 架构，可以看到支持很多 CPU 架构：
```bash
→ docker buildx ls

NAME/NODE    DRIVER/ENDPOINT             STATUS  PLATFORMS
mybuilder *  docker-container
  mybuilder0 unix:///var/run/docker.sock running linux/amd64, linux/arm64, linux/ppc64le, linux/s390x, linux/386, linux/arm/v7, linux/arm/v6
default      docker
  default    default                     running linux/amd64, linux/386
```

#### golang 最佳实践

**golang code**

```go
→ cat hello.go
package main

import (
        "fmt"
        "runtime"
)

func main() {
        fmt.Printf("Hello, %s!\n", runtime.GOARCH)
}
```

**dockerfile**

注意：golang 程序在多平台构建时，dockerfile 中不需要指定 GOOS，go build 会在 docker build 过程中基于 buildx 模拟的硬件作为 GOOS 进行编译；基础镜像必须是多架构镜像，如 golang:alpine 就是多架构镜像

```dockerfile
→ cat Dockerfile
FROM golang:alpine AS builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o hello .

FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/hello .
CMD ["./hello"]
```

**使用 buildx 构建 golang 程序的多平台架构镜像并推倒仓库**

```bash
→ docker buildx build -t registry/hello-multi-arch --platform=linux/arm,linux/arm64,linux/amd64 . --push
```

#### Java 最佳实践

对于Java来说，同样我们需要有一个支持包括ARM在内的多架构基础镜像，以常用的tomcat基础镜像为例，我们可以在dockerhub上搜索到官方的镜像。  
例如镜像`tomcat:jdk8-openjdk`，可以在tag页面确认已经支持`linux/amd64`和`linux/arm64`两种架构。
一个构建出war的Java项目，可以参考以下：  
**Dockerfile**

```dockerfile
FROM tomcat:jdk8-openjdk
ENV TZ=Asia/Shanghai LANG=C.UTF-8 LANGUAGE=C.UTF-8 LC_ALL=C.UTF-8
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /usr/local/tomcat
RUN rm -rf webapps/* \
    && mkdir webapps/ROOT
COPY target/*.war webapps/ROOT
RUN cd webapps/ROOT \
    && jar -xvf $(ls | grep '.war')

ENTRYPOINT ["catalina.sh", "run"]
```

如果是使用其他的JDK版本、Tomcat版本或者是使用Jar包部署，可自行修改以上Dockerfile。

**使用 buildx 构建 Java程序的多平台架构镜像并推到仓库**  
和 golang 一样，参考选择使用以上`docker manifests`或者`docker buildx`的方式，可以构建出多平台架构镜像并推到镜像仓库。

```bash
→ docker buildx build -t registry/java-multi-arch --platform=linux/arm64,linux/amd64 . --push
```
