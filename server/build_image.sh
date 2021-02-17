#!/usr/bin/env bash

cleanup() {
  rm -f Dockerfile
}
trap cleanup ERR EXIT

set -Eex

go generate

cat <<EOF >>Dockerfile
FROM golang:1.16
ENV GO111MODULE=on
RUN mkdir -p /app $GOPATH/src/github.com/talkiewalkie/talkiewalkie/
WORKDIR /app
COPY authenticated common migrations models repository unauthenticated $GOPATH/src/github.com/talkiewalkie/talkiewalkie/
COPY go.mod go.sum init.go server.go .env.prod /app
#RUN go mod download
RUN go get
RUN go build
CMD ["./talkiewalkie -env prod"]
EOF

docker build .
