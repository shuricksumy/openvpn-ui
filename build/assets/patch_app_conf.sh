#!/bin/bash

i=$((1 + $RANDOM % 1000))
echo "" >> /opt/openvpn-gui/conf/app.conf
echo "sessionname=beegosession_$i" >> /opt/openvpn-gui/conf/app.conf



