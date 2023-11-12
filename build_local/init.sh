cd ~

sudo apt-get update
wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
sudo tar -xvf go1.21.3.linux-amd64.tar.gz
sudo mv go /usr/local

export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
source ~/.profile

cd ~/openvpn-ui

go mod tidy
go mod vendor
go install github.com/beego/bee@latest

$GOPATH/bin/bee run -gendoc=true