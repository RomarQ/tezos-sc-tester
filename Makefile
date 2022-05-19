APP_NAME := visualtez-testing

BIN := api

VERSION := 0.0.7

ALL_PLATFORMS := linux/amd64 linux/arm64

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
DOCKER_REPO := "ghcr.io/romarq/visualtez-testing"

AMD64_IMAGE ?= alpine:3.15.4
ARM64_IMAGE ?= arm64v8/alpine:3.15.4

all: install build

install: download-tezos-client
	@go mod tidy
	@go mod vendor
	@go install github.com/swaggo/swag/cmd/swag@latest

test:
	@go test -cover -coverprofile=coverage.out -v ./...

download-tezos-client: .download-tezos-client
.download-tezos-client:
	@mkdir -p tezos-bin/amd64 tezos-bin/arm64
	@wget -O tezos-bin/amd64/tezos-client https://gitlab.com/tezos/tezos/-/jobs/2376802446/artifacts/raw/tezos-binaries/x86_64/tezos-client
	@wget -O tezos-bin/arm64/tezos-client https://gitlab.com/tezos/tezos/-/jobs/2376802447/artifacts/raw/tezos-binaries/arm64/tezos-client
	@chmod +x tezos-bin/amd64/tezos-client tezos-bin/arm64/tezos-client
	@touch .download-tezos-client

BUILD_DIRS := bin/$(OS)_$(ARCH)

$(BUILD_DIRS):
	@mkdir -p $@

all-build: $(addprefix build-, $(subst /,-, $(ALL_PLATFORMS)))

build-%:
	@$(MAKE) build                        \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst -, ,$*)) \
	    GOARCH=$(lastword $(subst -, ,$*))

build: install $(foreach bin, $(BIN), bin/$(OS)_$(ARCH)/$(bin).build)

bin/%.build: $(BUILD_DIRS)
	@sh -c "ARCH=$(ARCH) OS=$(OS) VERSION=$(VERSION) ./scripts/build.sh"
	@echo "Compilation complete: $</$(shell basename $*)"

build-docs:
	@${HOME}/go/bin/swag fmt .
	@${HOME}/go/bin/swag init -d cmd/api --parseDependency

start:
	@bin/$(OS)_$(ARCH)/$(BIN)

docker-build: build
	@docker build --tag $(DOCKER_REPO):$(VERSION)_amd64 -f Dockerfile.amd64 .
	@docker build --tag $(DOCKER_REPO):$(VERSION)_arm64 -f Dockerfile.arm64 .

docker-push: docker-build
	@docker push $(DOCKER_REPO):$(VERSION)_amd64
	@docker push $(DOCKER_REPO):$(VERSION)_arm64

version:
	@echo $(VERSION)

clean: bin-clean vendor-clean

vendor-clean:
	@rm -rf vendor

bin-clean:
	@rm -rf bin
