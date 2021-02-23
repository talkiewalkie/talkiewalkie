#!/usr/bin/env bash

cleanup() {
  rm -f Dockerfile
}
trap cleanup ERR EXIT

set -Eex

go generate

cat <<EOF >>Dockerfile
FROM golang:1.16 AS builder

WORKDIR /go/src/app
COPY .env.prod .
COPY . .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 go build -o talkiewalkie

FROM amd64/alpine

RUN apk update
RUN apk add --no-cache imagemagick

COPY --from=builder /go/src/app/migrations migrations
COPY --from=builder /go/src/app/.env.prod .env.prod
COPY --from=builder /go/src/app/.secrets .secrets
COPY --from=builder /go/src/app/talkiewalkie talkiewalkie
EXPOSE 8080

CMD ["./talkiewalkie", "-env", "prod"]
EOF

docker build --platform linux/amd64 -t gcr.io/talkiewalkie-305117/talkiewalkie-back:2 -t talkiewalkie-back .
