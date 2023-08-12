#!/bin/bash

set -e

time docker run \
    -v "$PWD/../":/go/src/github.com/shuricksumy/openvpn-ui \
    --rm \
    -w /usr/src/myapp \
    awalach/beego:1.8.1 \
    sh -c "cd /go/src/github.com/shuricksumy/openvpn-ui/ && bee version && bee pack -exr='^vendor|^data.db|^build|^README.md|^docs|^README_ORIGINAL.md|^screenshots|^pkg'"
