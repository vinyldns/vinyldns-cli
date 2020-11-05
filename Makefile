NAME=vinyldns
VERSION=0.9.1
TAG=v$(VERSION)
ARCH=$(shell uname -m)
PREFIX=/usr/local
DOCKER_NAME=vinyldns/vinyldns-cli
IMG=${DOCKER_NAME}:${VERSION}
LATEST=${DOCKER_NAME}:latest
VINYLDNS_REPO=github.com/vinyldns/vinyldns
VINYLDNS_VERSION=0.9.3
SRC=src/*.go
LOCAL_GO_PATH=`go env GOPATH`

all: test stop-api docker build-releases

install: build
	mkdir -p $(PREFIX)/bin
	cp -v bin/$(NAME) $(PREFIX)/bin/$(NAME)

uninstall:
	rm -vf $(PREFIX)/bin/$(NAME)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME) $(SRC)

build-releases:
	rm -rf release && mkdir release
	GOOS=linux go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_linux_$(ARCH) $(SRC)
	GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_darwin_$(ARCH) $(SRC)
	GOOS=linux CGO_ENABLED=0  go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_linux_$(ARCH)_nocgo $(SRC)

start-api:
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
	$(LOCAL_GO_PATH)/src/$(VINYLDNS_REPO)-$(VINYLDNS_VERSION)/bin/remove-vinyl-containers.sh

test-fmt:
	if [ `go fmt $(SRC) | wc -l` != "0" ]; then \
		echo "fix go code formatting by running 'go fmt'."; \
		exit 1; \
	fi;

test: test-fmt build start-api
	go get -u golang.org/x/lint/golint
	$(LOCAL_GO_PATH)/bin/golint -set_exit_status $(SRC)
	go vet $(SRC)
	go test $(SRC) -tags=integration -count=1

release: build-releases
	go get github.com/aktau/github-release
	github-release release \
		--user vinyldns \
		--repo vinyldns-cli \
		--tag $(TAG) \
		--name "$(TAG)" \
		--description "vinyldns-cli version $(VERSION)"
	cd release && ls | xargs -I FILE github-release upload \
		--user vinyldns \
		--repo vinyldns-cli \
		--tag $(TAG) \
		--name FILE \
		--file FILE

docker:
	docker build -t ${IMG} .
	docker tag ${IMG} ${LATEST}

docker-push:
	docker push ${LATEST}
	docker push ${IMG}

.PHONY: install uninstall build build_releases release test docker docker-push
