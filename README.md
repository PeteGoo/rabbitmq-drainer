# RabbitMQ Drainer
This is a simple app to drain queues in rabbit at a pre-defined rate

## Usage

### Drain a bunch of messages on a queue

```bash
rabbitmq-drainer drain \
  -host myrabbithost \
  -port 5672 \
  -user guest \
  -pass guest \
  -queue=hello \
  -concurrency 10 \
  -dwell 1
```
The above will drain 10 messages at a time waiting for 1ms after each message is drained.


## Testing
### Publish a bunch of messages

```
rabbitmq-drainer send -queue=hello -message="Hello World" -count 10000
```

### Receive a message or two

```bash
rabbitmq-drainer receive -queue=hello -count=1
```

## Building

### Building for Linux

```bash
docker run --rm -v "$PWD":/usr/src/rabbitmq-drainer -w /usr/src/rabbitmq-drainer -e GOMOD=/usr/src/rabbitmq-drainer/go.mod golang:1.13 go build -o ./out/linux/rabbitmq-drainer -v
```


### Building for Mac OS

```bash
docker run --rm -v "$PWD":/usr/src/rabbitmq-drainer -w /usr/src/rabbitmq-drainer -e GOMOD=/usr/src/rabbitmq-drainer/go.mod -e GOOS=darwin -e GOARCH=amd64 golang:1.13 go build -o ./out/darwin-amd64/rabbitmq-drainer -v
```


### Building for Windows

```bash
docker run --rm -v "$PWD":/usr/src/rabbitmq-drainer -w /usr/src/rabbitmq-drainer -e GOMOD=/usr/src/rabbitmq-drainer/go.mod -e GOOS=windows -e GOARCH=amd64 golang:1.13 go build -o ./out/win-amd64/rabbitmq-drainer.exe -v
```


