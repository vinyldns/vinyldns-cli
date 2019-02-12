NAME=vinyldns
VERSION=0.8.8
TAG=v$(VERSION)
ARCH=$(shell uname -m)
PREFIX=/usr/local
DOCKER_NAME=vinyldns/vinyldns-cli
IMG=${DOCKER_NAME}:${VERSION}
LATEST=${DOCKER_NAME}:latest
BATS=github.com/sstephenson/bats
VINYLDNS=github.com/vinyldns/vinyldns

all: lint vet acceptance stop-api build-releases

install: build
	mkdir -p $(PREFIX)/bin
	cp -v bin/$(NAME) $(PREFIX)/bin/$(NAME)

uninstall:
	rm -vf $(PREFIX)/bin/$(NAME)

build: deps
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME) src/*.go

build-releases: deps
	rm -rf release && mkdir release
	GOOS=linux  go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_linux_$(ARCH) src/*.go
	GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_darwin_$(ARCH) src/*.go
	GOOS=linux CGO_ENABLED=0  go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_linux_$(ARCH)_nocgo src/*.go

deps:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

start-api:
	if [ ! -d "$(GOPATH)/src/$(VINYLDNS)" ]; then \
		echo "$(VINYLDNS) not found in your GOPATH (necessary for acceptance tests), getting..."; \
		git clone https://$(VINYLDNS) $(GOPATH)/src/$(VINYLDNS); \
	fi
	$(GOPATH)/src/$(VINYLDNS)/bin/docker-up-vinyldns.sh \
		--api-only \
		--version 0.8.0

stop-api:
	$(GOPATH)/src/$(VINYLDNS)/bin/remove-vinyl-containers.sh

bats:
	if ! [ -x ${GOPATH}/src/${BATS}/bin/bats ]; then \
		git clone --depth 1 https://${BATS}.git ${GOPATH}/src/${BATS}; \
	fi

acceptance: build bats start-api
	${GOPATH}/src/${BATS}/bin/bats tests

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

lint: deps
	golint src/... -set_exit_status

vet:
	go tool vet src

docker:
	docker build -t ${IMG} .
	docker tag ${IMG} ${LATEST}

docker-push:
	docker push ${LATEST}
	docker push ${IMG}

.PHONY: install uninstall build build_releases deps release lint vet docker docker-push
