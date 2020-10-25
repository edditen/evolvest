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


.PHONY: gen-pb
gen-pb:
	protoc -I api/pb/evolvest/ api/pb/evolvest/evolvest.proto --go_out=plugins=grpc:api/pb/evolvest/


