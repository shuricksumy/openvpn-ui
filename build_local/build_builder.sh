#!/bin/bash

mkdir -p ./src/github.com/shuricksumy/openvpn-ui
cd ./src/github.com/shuricksumy/openvpn-ui
git clone https://github.com/shuricksumy/openvpn-ui.git ./

cd ./build

docker buildx create --use

./build_builder.sh