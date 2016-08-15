# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/workflow-cli

HOST_OS := $(shell uname)
ifeq ($(HOST_OS),Darwin)
	GOOS=darwin
else
	GOOS=linux
endif

# The latest git tag on branch
GIT_TAG := $(shell git describe --abbrev=0 --tags)
REVISION ?= $(shell git rev-parse --short HEAD)

REGISTRY ?= quay.io/
IMAGE_PREFIX ?= deisci
IMAGE := ${REGISTRY}${IMAGE_PREFIX}/workflow-cli-dev:${REVISION}

BUILD_OS ?=linux darwin windows
BUILD_ARCH ?=amd64 386

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.16.0
DEV_ENV_WORK_DIR := /go/src/${repo_path}
DEV_ENV_PREFIX := docker run --rm -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_PREFIX_CGO_ENABLED := docker run --rm -e CGO_ENABLED=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
DIST_DIR ?= _dist

GOTEST = go test --race

define check-static-binary
  if file $(1) | egrep -q "(statically linked|Mach-O)"; then \
    echo -n ""; \
  else \
    echo "The binary file $(1) is not statically linked. Build canceled"; \
    exit 1; \
  fi
endef

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

build: binary-build
	@$(call check-static-binary,deis)

# This is supposed to be run within a docker container
build-latest:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}-${REVISION}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/deis-latest-{{.OS}}-{{.Arch}}" .

# This is supposed to be run within a docker container
build-revision:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}-${REVISION}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/${REVISION}/deis-${REVISION}-{{.OS}}-{{.Arch}}" .

# This is supposed to be run within a docker container
build-tag:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/${GIT_TAG}/deis-${GIT_TAG}-{{.OS}}-{{.Arch}}" .

build-all: build-latest build-revision

binary-build:
	$(eval GO_LDFLAGS= -ldflags '-X ${repo_path}/version.Version=dev-${REVISION}')
	${DEV_ENV_PREFIX} -e GOOS=${GOOS} ${DEV_ENV_IMAGE} go build -a -installsuffix cgo ${GO_LDFLAGS} -o deis .

dist: build-all

install:
	cp deis $$GOPATH/bin

test-style: build-test-image
	docker run --rm ${IMAGE} lint

test: build-test-image
	docker run --rm ${IMAGE} test

build-test-image:
	docker build -t ${IMAGE} .

push-test-image: build-test-image
	docker push ${IMAGE}
