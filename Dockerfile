FROM golang:1.11.4-alpine3.8 as builder
WORKDIR /go/src/github.com/vinkdong/image-sync
ADD . .

RUN \
  apk add gcc build-base

RUN \
  go build .

FROM alpine

COPY --from=builder /go/src/github.com/vinkdong/image-sync/image-sync /image-sync

RUN \
set -ex \
   mkdir -p /etc/image-sync  && \
   apk add --no-cache ca-certificates

CMD ["/image-sync","sync","-c","/etc/image-sync/config.yml","-d"]