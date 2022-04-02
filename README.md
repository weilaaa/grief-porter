# Grief-Porter

people who move bricks should not yield

## Requirement

1. if you want to use `docker manifest` feature, must [open it](https://docs.docker.com/engine/reference/commandline/manifest/) at first
2. if you want to push image to private register, must `docker login` at first

## Quick Start

1. download specified binary from [release](https://github.com/weilaaa/grief-porter/releases)
2. `chmod +x ./goporter-*`
3. `mv ./goporter-* /usr/local/bin/goporter`
4. use it

```
goporter --config=./config.json
```

## Config

- auto move all images under given manifest
> *amend* means 'docker create manifest A B C --amend'

> *insecure* means add '--insecure' with 'docker manifest'
```json
[
  {
    "sources": [
      {
        "addr": "nginx:latest"
      }
    ],
    "amend": true,
    "auto": true,
    "insecure": true,
    "manifest": "registry.self/weilaaa/nginx:latest"
  }
]
```

- general move will lift given specified images into one manifest
```json
[
  {
    "sources": [
      {
        "addr": "golang:1.17@sha256:d860e175278037ee2429fecb1150bf10635ff4488c5a6faf695b169bf2c0868f",
        "new_tag": "registry.self/weilaaa/golang:1.17-amd64"
      },
      {
        "addr": "golang:1.17@sha256:6245b2bee36df7a76a983b7213af5765d6b61fda5a44fbaf95716135af152dac",
        "new_tag": "registry.self/weilaaa/golang:1.17-amd64-v8"
      }
    ],
    "auto": true,
    "manifest": "registry.self/weilaaa/golang:1.17-multi"
  }
]
```

- move images by specified platform rather than digest
```json
[
  {
    "sources": [
      {
        "addr": "golang:1.17",
        "platform": "amd64",
        "new_tag": "registry.self/weilaaa/golang:1.17-amd64"
      },
      {
        "addr": "golang:1.17",
        "platform": "arm64",
        "new_tag": "registry.self/weilaaa/golang:1.17-amd64"
      }
    ],
    "auto": true,
    "manifest": "registry.self/weilaaa/golang:1.17-multi"
  }
]
```

- move single image to another place without creating manifest
> you can remark whatever you want in config json file, but it effects nothing
```json
[
  {
    "sources": [
      {
        "addr": "golang:1.17@sha256:d860e175278037ee2429fecb1150bf10635ff4488c5a6faf695b169bf2c0868f",
        "remark": "my-image",
        "new_tag": "registry.self/weilaaa/golang:1.17-amd64"
      }
    ]
  }
]
```

## What is multi arch images

you can find tutorials [here](https://github.com/weilaaa/grief-porter/blob/main/docs/multi-arch-image.md)

## Local Usage

serially execute project operation

```bash
make run-serial
```

parallels execute project operation

```bash
make run-parallels
```