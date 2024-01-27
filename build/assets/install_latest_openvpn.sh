#!/bin/bash
.  "/etc/os-release"

arch=$(dpkg --print-architecture)
version=stable
osrelease=$VERSION_CODENAME

mkdir -p /etc/apt/keyrings
curl -fsSL https://swupdate.openvpn.net/repos/repo-public.gpg | gpg --dearmor > /etc/apt/keyrings/openvpn-repo-public.gpg
echo "deb [arch=${arch} signed-by=/etc/apt/keyrings/openvpn-repo-public.gpg] https://build.openvpn.net/debian/openvpn/${version} ${osrelease} main" > /etc/apt/sources.list.d/openvpn-aptrepo.list
apt update && apt install openvpn -y
# apt update && apt install openvpn-dco-dkms -y