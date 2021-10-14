FROM golang:1.16 AS builder

WORKDIR /go/src/app
COPY protos .
COPY server .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o talkiewalkie

FROM alpine

RUN apk update
RUN apk add --no-cache imagemagick

ARG GRPC_HEALTH_PROBE_VERSION=0.4.5
RUN wget -q -O /bin/grpc_health_probe "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64" && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /go/src/app/migrations migrations
COPY --from=builder /go/src/app/talkiewalkie talkiewalkie
EXPOSE 8080

CMD ["./talkiewalkie", "-env", "prod"]
