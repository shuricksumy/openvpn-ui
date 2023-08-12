#!/bin/bash

#export GOPATH=`pwd`

cd ./src/github.com/d3vilh/openvpn-ui

cd ./build

docker buildx create --use
#docker buildx create --name mycustombuilder --driver docker-container --bootstrap

./build_all.sh