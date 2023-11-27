#!/bin/bash

version="1.21.4"
export GOPATH=$HOME/go
export PATH="$HOME/go/bin:$PATH"

cd ~
wget https://go.dev/dl/go${version}.linux-amd64.tar.gz
chmod -R +w ${GOPATH}
rm -rf ${GOPATH} && tar -C ${HOME} -xzf go${version}.linux-amd64.tar.gz
rm ./go${version}.linux-amd64.tar.gz
go install github.com/beego/bee@latest