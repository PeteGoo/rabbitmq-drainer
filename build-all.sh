#!/bin/sh

docker run --rm -v "$PWD":/usr/src/rabbitmq-drainer -w /usr/src/rabbitmq-drainer -e GOMOD=/usr/src/rabbitmq-drainer/go.mod golang:1.13 go build -o ./out/linux/rabbitmq-drainer -v

docker run --rm -v "$PWD":/usr/src/rabbitmq-drainer -w /usr/src/rabbitmq-drainer -e GOMOD=/usr/src/rabbitmq-drainer/go.mod -e GOOS=darwin -e GOARCH=amd64 golang:1.13 go build -o ./out/darwin-amd64/rabbitmq-drainer -v

docker run --rm -v "$PWD":/usr/src/rabbitmq-drainer -w /usr/src/rabbitmq-drainer -e GOMOD=/usr/src/rabbitmq-drainer/go.mod -e GOOS=windows -e GOARCH=amd64 golang:1.13 go build -o ./out/win-amd64/rabbitmq-drainer.exe -v