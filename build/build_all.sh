#!/bin/bash

set -e

docker login
docker buildx build --platform linux/amd64,linux/arm64 -f Multi-arch.dockerfile -t shuricksumy/openvpn-ui:latest . --push --no-cache