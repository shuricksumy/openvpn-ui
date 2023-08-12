FROM debian:stable

RUN apt-get update && apt-get install golang-go git curl bzip2 -y

RUN mkdir -p /go/src/github.com/d3vilh/openvpn-ui
WORKDIR /go/src/github.com/d3vilh/openvpn-ui

# Uncomment for a multi-arch buildx of the main branch
RUN cd /go/src/github.com/d3vilh/openvpn-ui
RUN git clone https://github.com/shuricksumy/openvpn-ui ./

# Uncomment for a multi-arch buildx of the develop branch
# RUN git clone -b develop --single-branch https://github.com/bnhf/pivpn-tap-web-ui

RUN export GOPATH=/go/ && go install github.com/beego/bee@latest

