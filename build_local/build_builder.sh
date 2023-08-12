#!/bin/bash

#export GOPATH=`pwd`

cd ./src/github.com/shuricksumy/openvpn-ui

cd ./build

docker buildx create --use

./build_builder.sh