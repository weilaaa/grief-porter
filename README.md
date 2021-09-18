# Grief-Porter

people who move bricks should not yield

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

## Usage

serially execute project operation

```bash
make run-serial
```

parallels execute project operation

```bash
make run-parallels
```