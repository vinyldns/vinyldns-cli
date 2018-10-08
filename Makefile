NAME=vinyldns
VERSION=0.8.2
TAG=v$(VERSION)
ARCH=$(shell uname -m)
PREFIX=/usr/local
VETARGS?=-all

all: lint vet build_releases

install: build
	mkdir -p $(PREFIX)/bin
	cp -v bin/$(NAME) $(PREFIX)/bin/$(NAME)

uninstall:
	rm -vf $(PREFIX)/bin/$(NAME)

build: deps
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME)

build_releases: deps
	rm -rf release && mkdir release
	GOOS=linux  go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_linux_$(ARCH)
	GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o release/$(NAME)_$(VERSION)_darwin_$(ARCH)

deps:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go get github.com/golang/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

release: build_releases
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
	golint -set_exit_status

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet $(VETARGS)"
	@go tool vet $(VETARGS) . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: install uninstall build build_releases deps release lint vet
