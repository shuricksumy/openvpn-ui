#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Directory where OpenVPN configuration files are stored
OVDIR=/etc/openvpn

# Change to the /opt directory
cd /opt/

# Change to the OpenVPN GUI directory
cd /opt/openvpn-gui

# Create the database directory if it does not exist
mkdir -p db
echo "db dir created on this local path:"
pwd
echo "db dir contents:"
ls -lrt

# Set random session ID
if ! grep -qs "^sessionname=" /opt/openvpn-gui/conf/app.conf; then
    i=$((1 + $RANDOM % 1000))
    echo "" >> /opt/openvpn-gui/conf/app.conf
    echo "sessionname=beegosession_$i" >> /opt/openvpn-gui/conf/app.conf
fi

# Start the OpenVPN GUI
echo "Starting openvpn-ui!"
./openvpn-ui
