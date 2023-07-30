FROM debian:latest
WORKDIR /opt
EXPOSE 8080

RUN apt-get update && apt-get install -y curl

ADD assets/start.sh /opt/start.sh
ADD assets/generate_ca_and_server_certs.sh /opt/scripts/generate_ca_and_server_certs.sh
ADD assets/vars.template /opt/scripts/vars.template
ADD assets/genclient.sh /opt/scripts/genclient.sh
ADD assets/openvpn-install-v2.sh /opt/scripts/openvpn-install-v2.sh
ADD assets/restart.sh /opt/scripts/restart.sh
ADD assets/rmcert.sh /opt/scripts/rmcert.sh
ADD assets/rmclient.sh /opt/scripts/rmclient.sh


ADD openvpn-ui.tar.gz /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/data.db
ADD assets/app.conf /opt/openvpn-gui/conf/app.conf

CMD /opt/start.sh
