#!/bin/bash

if [ -z $1 ]; then 
    echo "ERROR: please set parameter - beta_branch_name !"
fi

# Multi-arch development build
# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f Multi-arch.dockerfile -t shurick/openvpn-ui:beta . #--push --no-cache
docker login
# docker buildx build --platform linux/amd64,linux/arm64 -f Multi-arch.dockerfile -t shuricksumy/openvpn-ui:beta . --push --no-cache

docker buildx build BRANCH=$1 --platform linux/amd64 -f Multi-arch-beta.dockerfile -t shuricksumy/openvpn-ui:beta . --push --no-cache


# Single-arch (amd64) development build
# docker buildx build --platform linux/amd64 -f Multi-arch.dockerfile -t shuricksumy/openvpn-ui . --push --no-cache

#rm -f $PKGFILE
