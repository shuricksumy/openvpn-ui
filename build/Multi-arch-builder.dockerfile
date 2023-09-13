FROM debian:stable

# Add build arguments for specific architectures
ARG TARGETPLATFORM


RUN apt-get update && apt-get install golang-go git curl bzip2 wget -y
RUN apt-get remove golang-go -y

RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz && tar -C /usr/local/ -xzf go1.21.1.linux-amd64.tar.gz; fi

RUN if [ "$TARGETPLATFORM" = "linux/arm64" ]; then wget https://go.dev/dl/go1.21.1.linux-arm64.tar.gz && tar -C /usr/local/ -xzf go1.21.1.linux-arm64.tar.gz; fi

ENV PATH="$PATH:/usr/local/go/bin"
RUN export PATH="$PATH:/usr/local/go/bin"; echo $PATH
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc

RUN mkdir -p /go/src/github.com/shuricksumy/openvpn-ui
WORKDIR /go/src/github.com/shuricksumy/openvpn-ui

RUN export GOPATH=/go/ && go install github.com/beego/bee@latest && /go/bin/bee update

