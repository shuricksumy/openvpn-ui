#!/bin/bash

set -e

#PKGFILE=openvpn-ui.tar.gz 

#cp -f ../$PKGFILE ./

# Multi-arch the manifest way -- uncomment for each architecture, save and run script
# docker build -t shuricksumy/openvpn-ui:manifest-amd64 --build-arg ARCH=amd64/ . --push --no-cache
# docker build -t shuricksumy/openvpn-ui:manifest-arm64 --build-arg ARCH=arm64/ . --push --no-cache
# docker build -t shuricksumy/openvpn-ui:manifest-armv7 --build-arg ARCH=armv7/ . --push --no-cache

# Multi-arch the buildx way -- just use the command below, don't run the script
# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f Multi-arch.dockerfile -t shuricksumy/openvpn-ui . --push --no-cache

# Multi-arch development build
# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f Multi-arch.dockerfile -t shurick/openvpn-ui:beta . #--push --no-cache
docker login
docker buildx build --platform linux/amd64,linux/arm64 -f Multi-arch.dockerfile -t shuricksumy/openvpn-ui:latest . --push --no-cache


# Single-arch (amd64) development build
# docker buildx build --platform linux/amd64 -f Multi-arch.dockerfile -t shuricksumy/openvpn-ui . --push --no-cache

#rm -f $PKGFILE
