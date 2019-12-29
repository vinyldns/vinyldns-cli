FROM golang:alpine AS builder

RUN apk add --update make git
COPY . /go/src/github.com/vinyldns/vinyldns-cli
RUN cd /go/src/github.com/vinyldns/vinyldns-cli \
  && make build-releases \
  && for i in $(head -n 2 Makefile); do eval $i; done \
  && cp /go/src/github.com/vinyldns/vinyldns-cli/release/${NAME}_${VERSION}_linux_$(uname -m)_nocgo /go/src/github.com/vinyldns/vinyldns-cli/vinyldns

FROM scratch

WORKDIR /root/
COPY --from=builder /go/src/github.com/vinyldns/vinyldns-cli/vinyldns .

ENTRYPOINT ["./vinyldns"]
