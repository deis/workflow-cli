# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/workflow-cli

HOST_OS := $(shell uname)
ifeq ($(HOST_OS),Darwin)
	GOOS=darwin
else
	GOOS=linux
endif

BUILD_OS ?=linux darwin windows
BUILD_ARCH ?=amd64 386

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.14.0
DEV_ENV_WORK_DIR := /go/src/${repo_path}
DEV_ENV_PREFIX := docker run --rm -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_PREFIX_CGO_ENABLED := docker run --rm -e CGO_ENABLED=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
DIST_DIR := _dist

GSUTIL_IMAGE := google/cloud-sdk:latest
GSUTIL_PREFIX := docker run --rm -v  ${CURDIR}/tmp:/.config -v ${CURDIR}/${DIST_DIR}:/upload
GSUTIL_CMD := ${GSUTIL_PREFIX} ${GSUTIL_IMAGE}
GCS_BUCKET ?= "gs://workflow-cli"

GO_FILES = $(wildcard *.go)
GO_LDFLAGS = -ldflags "-s -X ${repo_path}/version.BuildVersion=${VERSION}"
GO_PACKAGES = cmd parser $(wildcard pkg/*)
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))
GOFMT = gofmt -e -l -s
GOTEST = go test --cover --race -v

# The tag of the commit
GIT_TAG := $(shell git tag -l --contains HEAD)
VERSION ?= $(shell git rev-parse --short HEAD)

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
ifdef GIT_TAG
	${DEV_ENV_CMD} gox -verbose -parallel=3 ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="$(DIST_DIR)/${GIT_TAG}/deis-${GIT_TAG}-{{.OS}}-{{.Arch}}" .
else
	${DEV_ENV_CMD} gox -verbose -parallel=3 ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="$(DIST_DIR)/${VERSION}/deis-${VERSION}-{{.OS}}-{{.Arch}}" .
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

setup-gotools:
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/cover
	go get -u golang.org/x/tools/cmd/vet

test: test-style test-unit

test-style:
	${DEV_ENV_CMD} sh -c '${GOFMT} ${GO_FILES} ${GO_PACKAGES} && go vet $(repo_path) $(GO_PACKAGES_REPO_PATH)'
	@for i in $(addsuffix /...,$(GO_PACKAGES)); do \
		${DEV_ENV_CMD} golint $$i; \
	done

test-unit:
	${DEV_ENV_PREFIX_CGO_ENABLED} ${DEV_ENV_IMAGE} sh -c '${GOTEST} $$(glide nv)'

upload-gcs:
	${GSUTIL_CMD} sh -c 'gcloud auth activate-service-account -q --key-file /.config/key.json'
	${GSUTIL_CMD} sh -c 'gsutil -mq cp -a public-read -r /upload/* ${GCS_BUCKET}'
	# This has to run in the container to delete files created by the container
	${GSUTIL_CMD} sh -c 'rm -rf /.config/*'
