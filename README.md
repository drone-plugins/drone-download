# drone-download

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-download/status.svg)](http://cloud.drone.io/drone-plugins/drone-download)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/download.svg)](https://microbadger.com/images/plugins/download "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-download?status.svg)](http://godoc.org/github.com/drone-plugins/drone-download)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-download)](https://goreportcard.com/report/github.com/drone-plugins/drone-download)

Drone plugin to download files required for a build, also makes it possible to inject basic authentication in a secure way. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-download/).

## Build

Build the binary with the following command:

```bash
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-download ./cmd/drone-download
```

## Docker

Build the Docker image with the following command:

```bash
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/download .
```

## Usage

```bash
docker run --rm \
  -e PLUGIN_SOURCE=https://github.com/drone/drone-cli/releases/download/v0.8.5/drone_linux_amd64.tar.gz \
  -e PLUGIN_DESTINATION=drone_linux.tar.gz \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/download
```
