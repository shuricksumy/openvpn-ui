#!/bin/bash

export GOPATH=$HOME/go
export PATH="$HOME/go/bin:$PATH"
export SITE_NAME="TCP version"

export OPENVPN_ADMIN_USERNAME=admin # Leave this default as-is and update on first-run
export OPENVPN_ADMIN_PASSWORD=admin # Leave this default as-is and update on first-run
export OVDIR=/etc/openvpn

rm $GOPATH/bin/openvpn-ui
cd ../
#rm ./data.db
rm ./openvpn-ui

if [ ! -f ${OVDIR}/clientDetails.json ]; then
    touch ${OVDIR}/clientDetails.json
fi

go mod tidy
$GOPATH/bin/bee run -gendoc=false
