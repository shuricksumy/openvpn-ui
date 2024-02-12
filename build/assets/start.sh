#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Directory where OpenVPN configuration files are stored
if [[ $OVDIR == "" ]]; then
	OVDIR="/etc/openvpn"
fi

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

# Set site name
if ! grep -qs "^sitename=" /opt/openvpn-gui/conf/app.conf; then
    if [[ -n $SITE_NAME ]]; then
        echo "" >> /opt/openvpn-gui/conf/app.conf
        echo "sitename=${SITE_NAME}" >> /opt/openvpn-gui/conf/app.conf
    fi
else
    if [[ -n $SITE_NAME ]]; then
        sed -i '/sitename=/s/.*/sitename='"$SITE_NAME"'/' /opt/openvpn-gui/conf/app.conf
    fi
fi

# Set openvpn docker container name
if ! grep -qs "^OpenVpnServerDockerName=" /opt/openvpn-gui/conf/app.conf; then
    if [[ -n $OPENVPN_SERVER_DOCKER_NAME ]]; then
        echo "" >> /opt/openvpn-gui/conf/app.conf
        echo "OpenVpnServerDockerName=${OPENVPN_SERVER_DOCKER_NAME}" >> /opt/openvpn-gui/conf/app.conf
    fi
else
    if [[ -n $OPENVPN_SERVER_DOCKER_NAME ]]; then
        sed -i '/OpenVpnServerDockerName=/s/.*/OpenVpnServerDockerName='"$OPENVPN_SERVER_DOCKER_NAME"'/' /opt/openvpn-gui/conf/app.conf
    fi
fi

# Set openvpn management address
if ! grep -qs "^OpenVpnManagementAddress=" /opt/openvpn-gui/conf/app.conf; then
    if [[ -n $OPENVPN_MANAGEMENT_ADDRESS ]]; then
        echo "" >> /opt/openvpn-gui/conf/app.conf
        echo "OpenVpnManagementAddress=${OPENVPN_MANAGEMENT_ADDRESS}" >> /opt/openvpn-gui/conf/app.conf
    fi
else
    if [[ -n $OPENVPN_MANAGEMENT_ADDRESS ]]; then
        sed -i '/OpenVpnManagementAddress=/s/.*/OpenVpnManagementAddress='"$OPENVPN_MANAGEMENT_ADDRESS"'/' /opt/openvpn-gui/conf/app.conf
    fi
fi

# Set site name
if ! grep -qs "^httpport=" /opt/openvpn-gui/conf/app.conf; then
    if [[ -n $APP_PORT ]]; then
        echo "" >> /opt/openvpn-gui/conf/app.conf
        echo "httpport=${APP_PORT}" >> /opt/openvpn-gui/conf/app.conf
    fi
else
    if [[ -n $APP_PORT ]]; then
        sed -i '/httpport=/s/.*/httpport='"$APP_PORT"'/' /opt/openvpn-gui/conf/app.conf
    fi
fi

# Set URL PREFIX
if ! grep -qs "^BaseURLPrefix=" /opt/openvpn-gui/conf/app.conf; then
    if [[ -n $URL_PREFIX ]]; then
        echo "" >> /opt/openvpn-gui/conf/app.conf
        echo "BaseURLPrefix=${URL_PREFIX}" >> /opt/openvpn-gui/conf/app.conf
    fi
else
    if [[ -n $URL_PREFIX ]]; then
        sed -i '/BaseURLPrefix=/s/.*/BaseURLPrefix='"$URL_PREFIX"'/' /opt/openvpn-gui/conf/app.conf
    fi
fi


# if [ ! -f ${OVDIR}/clientDetails.json ]; then
#     touch ${OVDIR}/clientDetails.json
# fi

# if [ ! -f ${OVDIR}/routesDetails.json ]; then
#     touch ${OVDIR}/routesDetails.json
# fi

# wait for openvpn server ready
# until [ -f ${OVDIR}/server.conf ]
# do
#      sleep 1
# done

# Start the OpenVPN GUI
echo "Starting openvpn-ui!"
./openvpn-ui
