#!/usr/bin/env bash

cleanup() {
  rm -f Dockerfile
}
trap cleanup ERR EXIT

cat <<EOF >>Dockerfile
FROM golang:1.16 AS builder

WORKDIR /go/src/app
COPY protos .
COPY server .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o talkiewalkie

FROM alpine

RUN apk update
RUN apk add --no-cache imagemagick

COPY --from=builder /go/src/app/migrations migrations
COPY --from=builder /go/src/app/talkiewalkie talkiewalkie
EXPOSE 8080

CMD ["./talkiewalkie", "-env", "prod"]
EOF

set -Eex

sha="$(git rev-parse HEAD)"
docker build -t gcr.io/talkiewalkie-305117/talkiewalkie-back:${sha} -t talkiewalkie-back .
echo "built docker image 'gcr.io/talkiewalkie-305117/talkiewalkie-back:${sha}'"
