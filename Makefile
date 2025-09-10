.PHONY: test lint build docker-build run stop

test:
	go test -v -cover ./...

lint:
	golangci-lint run

build:
	go build -v -o bin/server ./cmd/server

docker-build:
	docker build -t mortgage-calculator .

run:
	./bin/server

stop:
	docker stop mortgage-calc && docker rm mortgage-calc