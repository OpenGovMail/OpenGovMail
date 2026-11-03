.PHONY: build run test clean docker-build

BINARY_NAME=autodiscover
DOCKER_IMAGE=opengovmail-autodiscover

build:
	go build -o $(BINARY_NAME) main.go

run:
	go run main.go

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE) .
