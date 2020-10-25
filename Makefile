.PHONY: all
all: fmt test

.PHONY: fmt
fmt:
	go mod tidy
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golint ./...




