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
# If the latest commit is tagged
TAGGED_COMMIT := $(shell git tag -l --contains HEAD)
REVISION ?= $(shell git rev-parse --short HEAD)

ifdef TAGGED_COMMIT
	VERSION ?= ${GIT_TAG}
else
	VERSION ?= ${GIT_TAG}-${REVISION}
endif

BUILD_OS ?=linux darwin windows
BUILD_ARCH ?=amd64 386

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.17.0
DEV_ENV_WORK_DIR := /go/src/${repo_path}
DEV_ENV_PREFIX := docker run --rm -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_PREFIX_CGO_ENABLED := docker run --rm -e CGO_ENABLED=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
DIST_DIR := _dist

GO_LDFLAGS = -ldflags "-s -X ${repo_path}/version.Version=${VERSION}"
GOTEST = go test --race

# UID and GID of local user
UID := $(shell id -u)
GID := $(shell id -g)

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

build-latest:
	${DEV_ENV_CMD} gox -verbose -parallel=3 ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="$(DIST_DIR)/deis-latest-{{.OS}}-{{.Arch}}" .

build-revision:
ifdef TAGGED_COMMIT
	${DEV_ENV_CMD} gox -verbose -parallel=3 ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="$(DIST_DIR)/${GIT_TAG}/deis-${GIT_TAG}-{{.OS}}-{{.Arch}}" .
else
	${DEV_ENV_CMD} gox -verbose -parallel=3 ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="$(DIST_DIR)/${REVISION}/deis-${REVISION}-{{.OS}}-{{.Arch}}" .
endif

build-all: build-latest build-revision

binary-build:
	${DEV_ENV_PREFIX} -e GOOS=${GOOS} ${DEV_ENV_IMAGE} go build -a -installsuffix cgo ${GO_LDFLAGS} -o deis .

dist: build-all

install:
	cp deis $$GOPATH/bin

installer: build
	@if [ ! -d makeself ]; then git clone -b single-binary https://github.com/deis/makeself.git; fi
	PATH=./makeself:$$PATH BINARY=deis makeself.sh --bzip2 --current --nox11 . \
		deis-cli-`cat deis-version`-`go env GOOS`-`go env GOARCH`.run \
		"Deis CLI" "echo \
		&& echo 'deis is in the current directory. Please' \
		&& echo 'move deis to a directory in your search PATH.' \
		&& echo \
		&& echo 'See http://docs.deis.io/ for documentation.' \
		&& echo"

test-style:
	${DEV_ENV_CMD} lint

test: test-style
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} sh -c '${GOTEST} $$(glide nv)'

test-cover: test-style
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} test-cover.sh

# Set local user as owner for files
fileperms:
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} chown -R ${UID}:${GID} .
