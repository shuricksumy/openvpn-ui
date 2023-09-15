#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

if [ -z "$1" ]; then
    echo "Use parameter: Client Name"
    exit 1
fi

echo 'Unrevoke certificate...'

OVPN_PATH=/etc/openvpn

keys_index_file="${OVPN_PATH}/easy-rsa/pki/index.txt"

# Copy easy-rsa variables
cd ${OVPN_PATH}/easy-rsa
linenumber="$(grep -n "/CN=$1"'$' $keys_index_file | cut -f1 -d: | head -1)"
fileline="$(grep -n "/CN=$1"'$' $keys_index_file | head -1)"
line="$(grep "/CN=$1"'$' $keys_index_file | head -1)"

columns_number="$(echo $line | awk -F' ' '{print NF;}')"
echo "Columns_number: $columns_number"

if [[ $columns_number -eq 6 ]] && [[ $line == R* ]]; then

	exp_time="$(echo $fileline | awk '{print $2}')"
	serial="$(echo $fileline | awk '{print $4}')"
	data="$(echo $fileline | awk '{print $5}')"
	details="$(echo $fileline | awk '{print $6}')"

	## TODO REMOVE R->V and revoke time
	sed -i "${linenumber}d" $keys_index_file
    echo -e "V\t$exp_time\t\t$serial\t$data\t$details" >> $keys_index_file
	
	## DO MAGIC
	cd ${OVPN_PATH}/easy-rsa/ || return

	EASYRSA_CRL_DAYS=3650 ./easyrsa gen-crl
	rm -f ${OVPN_PATH}/crl.pem
	cp ${OVPN_PATH}/easy-rsa/pki/crl.pem ${OVPN_PATH}/crl.pem
	chmod 644 ${OVPN_PATH}/crl.pem

	cp ${OVPN_PATH}/easy-rsa/pki/index.txt{,.bk}

	cp  ${OVPN_PATH}/easy-rsa/pki/revoked/certs_by_serial/$serial.crt   ${OVPN_PATH}/easy-rsa/pki/issued/$1.crt
    mv  ${OVPN_PATH}/easy-rsa/pki/revoked/certs_by_serial/$serial.crt   ${OVPN_PATH}/easy-rsa/pki/certs_by_serial/$serial.pem
    mv  ${OVPN_PATH}/easy-rsa/pki/revoked/private_by_serial/$serial.key ${OVPN_PATH}/easy-rsa/pki/private/$1.key
	mv  ${OVPN_PATH}/easy-rsa/pki/revoked/reqs_by_serial/$serial.req    ${OVPN_PATH}/easy-rsa/pki/reqs/$1.req

	echo ""
	echo "Certificate for client $1 unRevoked."

elif [[ $columns_number -eq 5 ]] && [[ $fileline == V* ]]; then

    echo "Certificate is already unrevoked and active"
    exit 0;

else
    echo "Error; Key index file may be corrupted."
    exit 1;

fi