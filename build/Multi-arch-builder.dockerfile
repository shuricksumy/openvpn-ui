FROM debian:stable

RUN apt-get update && apt-get install golang-go git curl bzip2 -y

RUN mkdir -p /go/src/github.com/shuricksumy/openvpn-ui
WORKDIR /go/src/github.com/shuricksumy/openvpn-ui

RUN export GOPATH=/go/ && go install github.com/beego/bee@latest && /go/bin/bee update

