
test:
	go test -v ./...

integration-test:
	go test -tags integration ./...

build:
	go build -o bin/telness-manager cmd/main.go

run:
	go run cmd/main.go

docker-build:
	docker build -t telness-manager_app .

up:
	docker-compose up

down:
	docker-compose down && docker rmi telness-manager_app