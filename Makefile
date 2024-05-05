GO=go

project=little-gpt/gpt-proxy

VERSION?=v0.0.1
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

os?=$(shell uname -s | tr '[A-Z]' '[a-z]')
arch?=amd64
cgo?=
cc?=


IMAGE_FLAGS :=
IMAGE_FLAGS := $(IMAGE_FLAGS) --build-arg baseBuilderImage=$(baseBuilderImage)
IMAGE_FLAGS := $(IMAGE_FLAGS) --build-arg baseImage=$(baseImage)
IMAGE_FLAGS := $(IMAGE_FLAGS) --build-arg arch=$(arch)

GIT_BUILD := `git rev-parse --short HEAD`
TARGET := gpt-proxy
DOCKER_IMAGE=little-gpt/gpt-proxy

BASENAME := $(shell basename ${PWD})

LDFLAGS += -X "$(project)/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(project)/version.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(project)/version.Version=$(VERSION)"
LDFLAGS += -X "$(project)/version.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

ifneq ($(os), )
export GOOS=$(os)
endif

ifneq ($(arch), )
BUILD := $(GIT_BUILD)-$(arch)
export GOARCH=$(arch)
endif

build: $(TARGET)

$(TARGET): $(SRC)
	@echo "go build -os $(os) -tags $@ -arch $(arch)"
	$(GO) build -ldflags '$(LDFLAGS)' $(project)/cmd/$@

image:
	docker build -f ./Dockerfile -t $(DOCKER_IMAGE):$(VERSION)-$(BUILD) .

clean:
	rm -f $(TARGET)