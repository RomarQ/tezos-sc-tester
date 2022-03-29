APP_NAME := visualtez-testing

BIN := api

VERSION := 0.0.1

ALL_PLATFORMS := linux/amd64 linux/arm linux/arm64 linux/ppc64le linux/s390x

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
BASEIMAGE ?= gcr.io/distroless/static

IMAGE := $(REGISTRY)/$(BIN)
TAG := $(VERSION)__$(OS)_$(ARCH)

BUILD_IMAGE ?= golang:1.18-alpine

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

# This will build the binary under ./.go and update the real binary if needed.
.PHONY: .go/%.stamp
.go/%.stamp: $(BUILD_DIRS)
	@echo "Making $</$(APP_NAME)-$(shell basename $*)"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/$<:/go/bin  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/sh -c "                                            \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        ./scripts/build.sh                                    \
	    "
	@if ! cmp -s .go/$* $*; then \
	    mv .go/$* $</$(APP_NAME)-$(shell basename $*);            \
	    date >$@;                              \
	fi

version:
	@echo $(VERSION)

$(BUILD_DIRS):
	@mkdir -p $@

clean: bin-clean

bin-clean:
	rm -rf .go bin
