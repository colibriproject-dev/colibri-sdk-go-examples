uid != id -u ${USER}

fmt:
	go fmt ./...

run:
	go run main.go

test: mock
	go-acc --covermode=set -o coverage.txt ./...
	grep -v -E "main.go|_mock.go" coverage.txt > filtered_coverage.txt
	mv filtered_coverage.txt coverage.txt

cover:
	go tool cover -html coverage.txt

mock:
	find . -type f -name "*_mock.go" -exec rm -f {} \;
	go generate -v ./...

build: doc
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s" -o application

update-module:
	go mod tidy

doc:
	swag init --pd --ot=json

# só dever ser executado uma única vez após realizar o clone do projeto
init: mock doc update-module
	cp .env_example .env

install-go-linux:
	sudo snap install go --channel=1.24/stable --classic
	echo 'export GOPATH="$$HOME/go"' >>~/.profile
	echo 'export PATH="$$PATH:$$GOPATH/bin"' >>~/.profile
	echo 'export PATH="$$PATH:$$GOPATH/bin:/usr/local/go/bin"' >>~/.profile
	. ~/.profile

install-dependencies:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/ory/go-acc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

migration-create: ## make migration-create name="new_migration_name"
	docker run --rm --user  $(uid) -v ./migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)

migration-up: ## make migration-up
	docker run --rm --user  $(uid) -v ./migrations:/migrations migrate/migrate -database postgres://postgres:postgres@localhost:5432/postgres -path /migrations up

migration-down: ## make migration-down
	docker run --rm --user  $(uid) -v ./migrations:/migrations migrate/migrate -database postgres://postgres:postgres@localhost:5432/postgres -path /migrations down
