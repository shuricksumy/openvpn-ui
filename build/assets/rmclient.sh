#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

OVDIR=/etc/openvpn

# .ovpn file path
DEST_FILE_PATH="$OVDIR/clients/$1.ovpn"

# Check if .ovpn file exists
if [[ ! -f $DEST_FILE_PATH ]]; then
    echo "User not found."
    exit 1
fi

# Fix index.txt by removing everything after pattern "/name=$1" in the line
sed -i'.bak' "s/\/name=${1}\/.*//" $OVDIR/easy-rsa/pki/index.txt

export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

echo 'Revoke certificate...'

# Copy easy-rsa variables
cd $OVDIR/easy-rsa
#cp /etc/openvpn/config/easy-rsa.vars ./vars

# Revoke certificate
./easyrsa revoke "$1"

echo 'Create new Create certificate revocation list (CRL)...'
./easyrsa gen-crl
chmod +r ./pki/crl.pem

echo 'Sync pki directory...'
#rm -rf /etc/openvpn/pki/*
#cp -r ./pki/. /etc/openvpn/pki
cp /etc/openvpn/easy-rsa/pki/crl.pem /etc/openvpn/crl.pem

echo 'Done!'
echo 'If you want to disconnect the user please restart the service using docker-compose restart openvpn.'
