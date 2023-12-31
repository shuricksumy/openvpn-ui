FROM debian:stable
WORKDIR /opt
EXPOSE 8080

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN apt-get update && apt-get install -y curl bzip2 jq git wget openvpn iptables openssl \
wget ca-certificates iproute2 sqlite3 procps

COPY ./build/assets/start.sh /opt/start.sh
COPY ./build/assets/vars.template /opt/scripts/vars.template
COPY ./build/assets/openvpn-install-v2.sh /opt/scripts/openvpn-install-v2.sh
COPY ./build/assets/install_pkg.sh /opt/scripts/install_pkg.sh
COPY ./build/assets/restart.sh /opt/scripts/restart.sh
COPY ./build/assets/rmcert.sh /opt/scripts/rmcert.sh
COPY ./build/assets/renew.sh /opt/scripts/renew.sh

RUN /opt/scripts/install_pkg.sh


RUN apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

COPY ./dist/openvpn-ui-${TARGETOS}-${TARGETARCH}${TARGETVARIANT}/openvpn-ui.tar.gz /opt/openvpn-gui/
RUN tar -zxf /opt/openvpn-gui/openvpn-ui.tar.gz --directory /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/openvpn-ui.tar.gz /opt/openvpn-gui/data.db
COPY ./build/assets/app.conf /opt/openvpn-gui/conf/app.conf

# Advise to open necassary ports
EXPOSE 1194/udp 8080/tcp

CMD /opt/start.sh