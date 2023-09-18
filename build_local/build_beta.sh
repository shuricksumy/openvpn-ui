#!/bin/bash

if [ -z $1 ]; then
    echo "ERROR: Need parametr - beta_branch_name !!!"
fi 



mkdir -p ./src/github.com/shuricksumy/openvpn-ui
cd ./src/github.com/shuricksumy/openvpn-ui
git clone https://github.com/shuricksumy/openvpn-ui.git ./
git checkout $1 


cd ./build

docker buildx create --use

./build_beta.sh $1