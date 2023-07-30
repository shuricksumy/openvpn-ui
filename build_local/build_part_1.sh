#!/bin/bash

export GOPATH=`pwd`

mkdir -p ./src/github.com/d3vilh/openvpn-ui
cd ./src/github.com/d3vilh/openvpn-ui
git clone https://github.com/shuricksumy/openvpn-ui.git ./

#git checkout use_new_path

go mod tidy
go install github.com/beego/bee@latest

$GOPATH/bin/bee run -gendoc=true
