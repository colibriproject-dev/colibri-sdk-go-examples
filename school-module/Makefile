fmt:
	go fmt ./...

mock:
	go generate -v ./...

test: mock
	go test ./... --coverprofile coverage.out

cover:
	go tool cover -html coverage.out

run:
	go run main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s" -o application

update-module:
	go mod tidy

init: mock
	go mod tidy
	cp .env_example .env

install-go:
	sudo snap install go --channel=1.18/stable --classic

install: install-go
	go install github.com/golang/mock/mockgen@v1.6.0
	echo 'export GOPATH="$$HOME/go"' >>~/.profile
	echo 'export PATH="$$PATH:$$GOPATH/bin"' >>~/.profile
	echo 'export PATH="$$PATH:$$GOPATH/bin:/usr/local/go/bin"' >>~/.profile
	. ~/.profile
