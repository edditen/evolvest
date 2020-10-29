SHELL=/bin/bash
# build program
SRC_PATH = cmd
CFG_FILE = conf/config.yaml
BUILD_PATH = bin
SRC_FILES := $(shell cd $(SRC_PATH); find . -maxdepth 1 -type d|grep -v '/common')
TAGET := $(basename $(patsubst ./%,%,$(SRC_FILES)))
FILES := $(basename $(patsubst ./%,$(BUILD_PATH)/%,$(SRC_FILES)))

.PHONY: build $(TAGET) test-integration

# Example:
#   make build
#   make build GOFLAGS=-race
build:$(FILES)

$(TAGET):
	make -B $(BUILD_PATH)/$@

$(BUILD_PATH)/%: $(SRC_PATH)/%
	./scripts/build_script.sh $@ ./$^

.PHONY: clean
clean:
	rm -f $(FILES)

.PHONY: run
run:
	$(BUILD_PATH)/evolvestd -v -c $(CFG_FILE)


.PHONY: fmt
fmt:
	go mod tidy
	go fmt ./...
	golint ./...

.PHONY: test
test:
	go test -short -v ./...

.PHONY: gen-pb
gen-pb:
	protoc -I api/pb/evolvest/ api/pb/evolvest/evolvest.proto --go_out=plugins=grpc:api/pb/evolvest/


.PHONY: all
all: clean test build run