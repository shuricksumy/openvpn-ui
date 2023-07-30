#!/bin/bash

set -e

PKGFILE=openvpn-ui.tar.gz

cp -f ../$PKGFILE ./

docker build -t shuricksumy/openvpn-ui . --push --no-cache

rm -f $PKGFILE
