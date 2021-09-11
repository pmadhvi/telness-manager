
test:
	go test -v ./...

integration-test:
	go test -tags integration ./...

build:
	go build -o bin/iban-validator cmd/main.go

run:
	go run cmd/main.go

build-linux-64:
	GOOS=linux GOARCH=amd64 go build -o bin/iban-validator-linux-64 cmd/main.go

build-linux-32:
	GOOS=linux GOARCH=386 go build -o bin/iban-validator-linux-32 cmd/main.go


docker-build:
	docker build -t iban-validator .

up:
	docker-compose up

down:
	docker-compose down