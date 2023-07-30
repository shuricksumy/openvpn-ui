#!/bin/bash

export GOPATH=`pwd`

cd ./src/github.com/d3vilh/openvpn-ui

$GOPATH/bin/bee pack -exr='^vendor|^data.db|^build|^README.md|^docs'

