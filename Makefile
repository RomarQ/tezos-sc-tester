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

build-%:
	@$(MAKE) build                        \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

all-build: $(addprefix build-, $(subst /,_, $(ALL_PLATFORMS)))

build: $(foreach bin, $(BIN), bin/$(OS)_$(ARCH)/$(bin))

BUILD_DIRS := bin/$(OS)_$(ARCH)     \
              .go/bin/$(OS)_$(ARCH) \
              .go/cache

bin/%: .go/bin/%.stamp
	@true

.PHONY: .go/%.stamp
.go/%.stamp: $(BUILD_DIRS)
	@sh -c "ARCH=$(ARCH) OS=$(OS) VERSION=$(VERSION) ./scripts/build.sh"
	@echo "Compilation complete: $</$(APP_NAME)-$(shell basename $*)"

version:
	@echo $(VERSION)

$(BUILD_DIRS):
	@mkdir -p $@

clean: bin-clean vendor-clean

vendor-clean:
	rm -rf vendor

bin-clean:
	rm -rf .go bin
