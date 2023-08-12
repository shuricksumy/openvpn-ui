#!/bin/bash

cd ./src/github.com/shuricksumy/openvpn-ui

cd ./build

docker buildx create --use

./build_all.sh