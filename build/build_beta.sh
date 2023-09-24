#!/bin/bash

if [ -z $1 ]; then 
    echo "ERROR: please set parameter - beta_branch_name !"
fi

docker login
docker buildx build --build-arg BRANCH=$1 --platform linux/amd64 -f Multi-arch-beta.dockerfile -t shuricksumy/openvpn-ui:beta . --push --no-cache
