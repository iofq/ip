all: lint test install

lint:
	go vet ./...

test:
	go test ./...

install:
	go install ./...
