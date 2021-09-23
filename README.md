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

```json
[
  {
    "sources": [
      {
        "addr": "prom/prometheus:v2.18.2@sha256:cba0deaa490dea181e59df5ce8c10a0eb2c1aa0196f26c7eaade947448ae393a",
        "remark": "linux/arm64",
        "new_tag": "my.register.com/foo/prometheus:v2.18.2-arm64"
      },
      {
        "addr": "prom/prometheus:v2.18.2@sha256:9564f635c7d83bd242589842741bac3cf2746e9f94c250384850cf18ae09999d",
        "new_tag": "my.register.com/foo/prometheus:v2.18.2-amd64"
      }
    ],
    "manifest": "my.register.com/foo/prometheus:v2.18.2-multi"
  }
]
```

## Local Usage

serially execute project operation

```bash
make run-serial
```

parallels execute project operation

```bash
make run-parallels
```