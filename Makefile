APP_NAME := visualtez-testing

BIN := api

VERSION := 0.0.1

ALL_PLATFORMS := linux/amd64 linux/arm64

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

all: install build

install:
	@go mod tidy
	@go mod vendor

test:
	@go test -cover -coverprofile=coverage.out -v ./...

BUILD_DIRS := bin/$(OS)_$(ARCH)

$(BUILD_DIRS):
	@mkdir -p $@

all-build: $(addprefix build-, $(subst /,-, $(ALL_PLATFORMS)))

build-%:
	@$(MAKE) build                        \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst -, ,$*)) \
	    GOARCH=$(lastword $(subst -, ,$*))

build: $(foreach bin, $(BIN), bin/$(OS)_$(ARCH)/$(bin).build)

bin/%.build: $(BUILD_DIRS)
	@sh -c "ARCH=$(ARCH) OS=$(OS) VERSION=$(VERSION) ./scripts/build.sh"
	@echo "Compilation complete: $</$(shell basename $*)"

start:
	@bin/$(OS)_$(ARCH)/$(BIN)

version:
	@echo $(VERSION)

clean: bin-clean vendor-clean

vendor-clean:
	@rm -rf vendor

bin-clean:
	@rm -rf bin
