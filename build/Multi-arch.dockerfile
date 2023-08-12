FROM shuricksumy/builder:latest

#RUN apt-get update && apt-get install golang-go git curl bzip2 -y

#RUN mkdir -p /go/src/github.com/d3vilh/openvpn-ui
#WORKDIR /go/src/github.com/d3vilh/openvpn-ui

# Uncomment for a multi-arch buildx of the main branch
#RUN cd /go/src/github.com/d3vilh/openvpn-ui
#RUN git clone https://github.com/shuricksumy/openvpn-ui ./

# Uncomment for a multi-arch buildx of the develop branch
# RUN git clone -b develop --single-branch https://github.com/bnhf/pivpn-tap-web-ui

#RUN export GOPATH=/go/ && go install github.com/beego/bee@latest

RUN cd /go/src/github.com/d3vilh/openvpn-ui
WORKDIR /go/src/github.com/d3vilh/openvpn-ui

#RUN go mod tidy && /root/go/bin/bee pack -mod=readonly -exr='^vendor|^data.db|^build|^README.md|^docs|^README_ORIGINAL.md|^screenshots'

RUN export GOPATH=/go/ && go mod tidy && /go/bin/bee pack -exr='^vendor|^data.db|^build|^README.md|^docs|^README_ORIGINAL.md|^screenshots'
#-mod=readonly

FROM debian:stable
WORKDIR /opt
EXPOSE 8080


RUN apt-get update && apt-get install -y curl bzip2

COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/start.sh /opt/start.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/generate_ca_and_server_certs.sh /opt/scripts/generate_ca_and_server_certs.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/vars.template /opt/scripts/vars.template
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/genclient.sh /opt/scripts/genclient.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/openvpn-install-v2.sh /opt/scripts/openvpn-install-v2.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/install_pkg.sh /opt/scripts/install_pkg.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/restart.sh /opt/scripts/restart.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/rmcert.sh /opt/scripts/rmcert.sh
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/rmclient.sh /opt/scripts/rmclient.sh

RUN /opt/scripts/install_pkg.sh

COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/openvpn-ui.tar.gz /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/data.db
COPY --from=0  /go/src/github.com/d3vilh/openvpn-ui/build/assets/app.conf /opt/openvpn-gui/conf/app.conf


CMD /opt/start.sh