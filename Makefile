all: deps
	export GOPATH=`pwd`; go build src/communitty.go

clean:
	rm -f communitty
	rm -rf pkg
	rm -rf src/code.google.com
	rm -rf src/github.com

deps:
	export GOPATH=`pwd`; go get github.com/kr/pty
	export GOPATH=`pwd`; go get code.google.com/p/go.net/websocket
