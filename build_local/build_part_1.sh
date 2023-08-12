#!/bin/bash

export GOPATH=`pwd`

mkdir -p ./src/github.com/shuricksumy/openvpn-ui
cd ./src/github.com/shuricksumy/openvpn-ui
git clone https://github.com/shuricksumy/openvpn-ui.git ./

#git checkout use_new_path

go mod tidy
go mod vendor
go install github.com/beego/bee@latest

$GOPATH/bin/bee run -gendoc=true
