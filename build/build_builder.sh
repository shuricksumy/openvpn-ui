#!/bin/bash

set -e

# Multi-arch development build
docker login --username shuricksumy
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f Multi-arch-builder.dockerfile -t shuricksumy/builder:latest . --push --no-cache