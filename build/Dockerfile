FROM debian:stable
WORKDIR /opt
EXPOSE 8080

RUN apt-get update && apt-get install -y curl bzip2 jq openvpn iptables openssl wget ca-certificates iproute2

ADD assets/start.sh /opt/start.sh
ADD assets/vars.template /opt/scripts/vars.template
ADD assets/openvpn-install-v2.sh /opt/scripts/openvpn-install-v2.sh
ADD assets/EasyRSA-3.1.7.tgz ~/EasyRSA-3.1.7.tgz
ADD assets/install_pkg.sh /opt/scripts/install_pkg.sh
ADD assets/restart.sh /opt/scripts/restart.sh
ADD assets/rmcert.sh /opt/scripts/rmcert.sh
ADD assets/renew.sh /opt/scripts/renew.sh


RUN /opt/scripts/install_pkg.sh

ADD openvpn-ui.tar.gz /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/data.db
ADD assets/app.conf /opt/openvpn-gui/conf/app.conf

CMD /opt/start.sh
