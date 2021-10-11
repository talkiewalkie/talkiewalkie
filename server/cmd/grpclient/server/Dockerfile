FROM golang:1.16 AS builder

WORKDIR /go/src/app
COPY . .
COPY hwpb ./hwpb

RUN ls -lah .
RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o talkiewalkie
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o talkiewalkie

FROM alpine

COPY --from=builder /go/src/app/talkiewalkie talkiewalkie
EXPOSE 50051

CMD ["./talkiewalkie"]
