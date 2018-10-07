FROM golang:latest AS builder

RUN mkdir -p /go/src/github.com/vinyldns/vinyldns-cli \
    && git clone -b docker-file --single-branch \
        https://github.com/marekfilip/vinyldns-cli.git \
        /go/src/github.com/vinyldns/vinyldns-cli \
    && cd /go/src/github.com/vinyldns/vinyldns-cli \
    && make \
    && for i in $(head -n 2 Makefile); do eval $i; done \
    && cp /go/src/github.com/vinyldns/vinyldns-cli/release/${NAME}_${VERSION}_linux_$(uname -m)_nocgo /go/src/github.com/vinyldns/vinyldns-cli/vinyldns

FROM scratch

WORKDIR /root/
COPY --from=builder /go/src/github.com/vinyldns/vinyldns-cli/vinyldns .

ENTRYPOINT ["./vinyldns"]
