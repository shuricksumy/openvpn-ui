#!/bin/bash

export GOPATH=`pwd`

cd ./src/github.com/shuricksumy/openvpn-ui

$GOPATH/bin/bee pack -exr='^vendor|^data.db|^build|^README.md|^docs|^README_ORIGINAL.md|^screenshots|^pkg'

