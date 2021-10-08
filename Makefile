VERSION=0.10.0

SHELL=bash
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

NAME=vinyldns
TAG=v$(VERSION)

ARCH=$(shell uname -m)
ARCH_ARM=arm64
INSTALL_PATH=/usr/local
DOCKER_NAME=vinyldns/vinyldns-cli
IMG=${DOCKER_NAME}:${VERSION}
LATEST=${DOCKER_NAME}:latest

BATS=github.com/sstephenson/bats
VINYLDNS_REPO=github.com/vinyldns/vinyldns
VINYLDNS_VERSION=0.9.10

SOURCE_PATH:=$(ROOT_DIR)/src
LOCAL_GO_PATH=`go env GOPATH`

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
PLATFORMS=darwin linux windows

# Check that the required version of make is being used
REQ_MAKE_VER:=3.82
ifneq ($(REQ_MAKE_VER),$(firstword $(sort $(MAKE_VERSION) $(REQ_MAKE_VER))))
   $(error The version of MAKE $(REQ_MAKE_VER) or higher is required; you are running $(MAKE_VERSION))
endif

.ONESHELL:


.PHONY: install uninstall build build_releases release test docker docker-push

all: test build-releases

install: build
	@set -euo pipefail
	mkdir -p $(INSTALL_PATH)/bin
	cp -v bin/$(NAME) $(INSTALL_PATH)/bin/$(NAME)

uninstall:
	@set -euo pipefail
	rm -vf $(INSTALL_PATH)/bin/$(NAME)

build:
	@set -euo pipefail
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME) $(SOURCE_PATH)

build-releases:
	@set -euo pipefail
	rm -rf release && mkdir release
	for platform in $(PLATFORMS); do
	    GOOS=$${platform}
	    GOARCH=amd64
	    BINARY="$(NAME)"
	    if [ "$${platform}" == "windows" ]; then BINARY="$${BINARY}.exe"; fi

	    echo -n "Building $${BINARY} v$(VERSION) for $${platform}/$${GOARCH}..."
	    GOOS=$${platform} GOARCH=$${GOARCH} go build -ldflags "-X main.version=$(VERSION)" -o $(ROOT_DIR)/release/$${GOOS}_$${GOARCH}/$${BINARY} $(SOURCE_PATH);
	    echo -n "compressing..."
	    tar czf $(ROOT_DIR)/release/$(NAME)_$(VERSION)_$${GOOS}_$${GOARCH}.tar.gz -C $(ROOT_DIR)/release/$${GOOS}_$${GOARCH} $${BINARY};
	    echo "done."
	done

start-api:
	@set -euo pipefail
	if [ ! -d "$(LOCAL_GO_PATH)/src/$(VINYLDNS_REPO)-$(VINYLDNS_VERSION)" ]; then \
		echo "$(VINYLDNS_REPO)-$(VINYLDNS_VERSION) not found in your GOPATH (necessary for acceptance tests), getting..."; \
		git clone \
			--branch v$(VINYLDNS_VERSION) \
			https://$(VINYLDNS_REPO) \
		$(LOCAL_GO_PATH)/src/$(VINYLDNS_REPO)-$(VINYLDNS_VERSION); \
	fi
	$(LOCAL_GO_PATH)/src/$(VINYLDNS_REPO)-$(VINYLDNS_VERSION)/bin/docker-up-vinyldns.sh \
		--api-only \
		--version $(VINYLDNS_VERSION)

stop-api:
	@set -euo pipefail
	$(LOCAL_GO_PATH)/src/$(VINYLDNS_REPO)-$(VINYLDNS_VERSION)/bin/remove-vinyl-containers.sh

bats:
	@set -euo pipefail
	if ! [ -x ${LOCAL_GO_PATH}/src/${BATS}/bin/bats ]; then
		git clone --depth 1 https://${BATS}.git ${LOCAL_GO_PATH}/src/${BATS};
	fi

test-fmt:
	@set -euo pipefail
	if [ `go fmt $(SOURCE_PATH) | wc -l` != "0" ]; then
		echo "Fix go code formatting by running 'make format'."
		exit 1
	fi;

format:
	@set -euo pipefail
	go fmt $(SOURCE_PATH)

test: test-fmt build bats start-api
	@set -euo pipefail
	trap 'make stop-api' TERM INT EXIT
	go install golang.org/x/lint/golint@latest
	$(LOCAL_GO_PATH)/bin/golint -set_exit_status $(SOURCE_PATH)
	go vet $(SOURCE_PATH)
	${LOCAL_GO_PATH}/src/${BATS}/bin/bats tests

docker:
	@set -euo pipefail
	docker build -t ${IMG} .
	docker tag ${IMG} ${LATEST}

docker-push:
	@set -euo pipefail
	docker push ${LATEST}
	docker push ${IMG}


