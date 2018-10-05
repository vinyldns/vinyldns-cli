FROM golang:latest AS builder

RUN go get github.com/vinyldns/vinyldns-cli \
    && cd /go/src/github.com/vinyldns/vinyldns-cli \
    && go get github.com/golang/lint/golint \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o vinyldns .

FROM scratch

WORKDIR /root/
COPY --from=builder /go/src/github.com/vinyldns/vinyldns-cli/vinyldns .

ENTRYPOINT ["./vinyldns"]