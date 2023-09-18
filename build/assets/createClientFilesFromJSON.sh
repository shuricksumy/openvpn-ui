#!/bin/bash
# requred jq - apt install jq


# converts IPv4 as "A.B.C.D" to integer
# ip4_to_int 192.168.0.1
# => 3232235521
ip4_to_int() {
  IFS=. read -r i j k l <<EOF
$1
EOF
  echo $(( (i << 24) + (j << 16) + (k << 8) + l ))
}

# converts interger to IPv4 as "A.B.C.D"
#
# int_to_ip4 3232235521
# => 192.168.0.1
int_to_ip4() {
  echo "$(( ($1 >> 24) % 256 )).$(( ($1 >> 16) % 256 )).$(( ($1 >> 8) % 256 )).$(( $1 % 256 ))"
}

get_next_ip() {
    old=$(ip4_to_int $1)
    old=$((old+1))
    new=$(int_to_ip4 $old)
    echo $new
}

# validate IPv4 as "A.B.C.D"
#
# valid_ipv4 192.168.0.1
# => true/flase
valid_ipv4() {
echo $1 | grep -E -o "(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\
\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?$)" && return 
}

# returns the ip part of an CIDR
#
# cidr_ip "172.16.0.10/22"
# => 172.16.0.10
cidr_ip() {
  IFS=/ read -r ip _ <<EOF
$1
EOF
  echo $ip
}

# returns the prefix part of an CIDR
#
# cidr_prefix "172.16.0.10/22"
# => 22
cidr_prefix() {
  IFS=/ read -r _ prefix <<EOF
$1
EOF
  echo $prefix
}

# returns net mask in numberic from prefix size
#
# netmask_of_prefix 8
# => 4278190080
netmask_of_prefix() {
  echo $((4294967295 ^ (1 << (32 - $1)) - 1))
}

# returns default gateway address (network address + 1) from CIDR
cidr_default_gw() {
  ip=$(ip4_to_int $(cidr_ip $1))
  prefix=$(cidr_prefix $1)
  netmask=$(netmask_of_prefix $prefix)
  gw=$((ip & netmask + 1))
  int_to_ip4 $gw
}

# returns default gateway address (broadcast address - 1) from CIDR
cidr_default_gw_2() {
  ip=$(ip4_to_int $(cidr_ip $1))
  prefix=$(cidr_prefix $1)
  netmask=$(netmask_of_prefix $prefix)
  broadcast=$(((4294967295 - netmask) | ip))
  int_to_ip4 $((broadcast - 1))
}

# ======================================================================

OVDIR="/etc/openvpn"
USERDiR="${OVDIR}/ccd"
JSON="clientDetails.json"

if [ -z $JSON ]; then
    exit 1
fi

TIMESTAMP=$(date +%F_%T)

IFS=','
Clients=$(cat ${OVDIR}/${JSON} | jq -c 'map(.ClientName)' | sed 's/[][]//g'| sed 's/["]//g')

for client in ${Clients[*]}; do
    fileToSave=${USERDiR}/${client}
    echo "Client: $client"
    clientDetails=$(cat ${OVDIR}/${JSON} | jq -c ".[] | select( any(.; .ClientName == \"$client\") )")

    IFS=','
    ClientSelectedRoutes=$(echo "$clientDetails" | jq -c ".RouteListSelected" | sed 's/[][]//g' | sed 's/["]//g')

    
    ####
    echo "# Automatic generated client settings file - $TIMESTAMP" > ${fileToSave}
    ####

    staticIP=$(echo "$clientDetails" | jq -c ".StaticIP" | sed 's/["]//g')
    if [[ $(valid_ipv4 "$staticIP") ]]; then 
        nextIP=$(get_next_ip $staticIP)

        ####
        echo "ifconfig-push $staticIP $nextIP" >> ${fileToSave}
        ####
    fi

    isRouter=$(echo "$clientDetails" | jq -c ".IsRouter")
    if [[ $isRouter == "true" ]]; then

        thisRouterSubnet=$(echo "$clientDetails" | jq -c ".RouterSubnet" | sed 's/["]//g')
        thisRouterMask=$(echo "$clientDetails" | jq -c ".RouterMask" | sed 's/["]//g')

        isParamOK=true

        if [ $thisRouterMask == "" || $thisRouterSubnet == "" ]; then
            echo ">>>>> Network/Mask is empty"
            isParamOK=false
        fi

        if [ ! $(valid_ipv4 "$thisRouterSubnet") ]; then 
            echo ">>>>> ThisRouterSubnet IP is invalid: $thisRouterSubnet"
            isParamOK=false
        fi
        if [ ! $(valid_ipv4 "$thisRouterMask") ]; then 
            echo ">>>>> ThisRouterMask IP is invalid: $thisRouterMask"
            isParamOK=false
        fi

        if [ $isParamOK == true ]; then
            ####
            echo "iroute ${thisRouterSubnet} ${thisRouterMask}" >> ${fileToSave}
            ####
        fi
    fi
    ####
    echo "  " >> ${fileToSave}
    ####

    isRouteDefault=$(echo "$clientDetails" | jq -c ".IsRouteDefault")
    if [[ $isRouteDefault == "true" ]]; then
        ####
        echo "# Set VPN as default route" >> ${fileToSave}
        echo "push \"redirect-gateway def1\"" >> ${fileToSave}
        ####
    else
        ####
        echo "# Set VPN as default route" >> ${fileToSave}
        echo "# push \"redirect-gateway def1\"" >> ${fileToSave}
        ####
    fi
    ####
    echo "  " >> ${fileToSave}
    ####

    for route in ${ClientSelectedRoutes[*]}; do
        routeClient=$(cat ${OVDIR}/${JSON} | jq -c ".[] | select( any(.; .ClientName == \"$route\") )")
        routerSubnet=$(echo "$routeClient" | jq -c ".RouterSubnet" | sed 's/["]//g')

        if [ ! $(valid_ipv4 "$routerSubnet") ]; then 
            echo ">>>>> RouterSubnet IP is invalid: $routerSubnet"
            continue
        fi
        routerMask=$(echo "$routeClient" | jq -c ".RouterMask" | sed 's/["]//g')
        if [ ! $(valid_ipv4 "$routerMask") ]; then 
            echo ">>>>> RouterMask IP is invalid: $routerMask"
            continue
        fi
        description=$(echo "$routeClient" | jq -c ".Description" | sed 's/["]//g')
        
        ####
        echo "# Route to ${route} [${description}] device internal subnet" >> ${fileToSave}
        echo "push \"route ${routerSubnet} ${routerMask}\"" >> ${fileToSave}
        ####
    done
done
