.PHONY: run build test clean docker

run:
	go run ./cmd/server

build:
	go build -o bin/shortly ./cmd/server

test:
	go test ./... -v -cover

clean:
	rm -rf bin/

docker:
	docker-compose up --build -d

lint:
	golangci-lint run ./...
