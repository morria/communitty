all:
	export GOPATH=`pwd`; go build src/communitty.go

clean:
	rm -f communitty
	rm -rf pkg
	rm -rf src/golang.org
	rm -rf src/github.com

deps:
	export GOPATH=`pwd`; go get github.com/kr/pty
	export GOPATH=`pwd`; go get golang.org/x/net/websocket
