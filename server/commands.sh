#!/usr/bin/env bash

start_postgres() {
  pg_ctl -D /opt/homebrew/var/postgresql@12 start
}

install_golang_migrate_cli() {
  brew install golang-migrate
}

nukedb() {
  dropdb talkiewalkie && createdb talkiewalkie
  migrate -path migrations -database postgres://theo:pinguy@localhost:5432/talkiewalkie?sslmode=disable up
}

new_migration() {
  migrate create -ext sql -dir migrations/ -seq -digits 3 "$1"
}

new_secret() {
  openssl rand -hex 32
}

install_sqlboiler_cli() {
  go get github.com/volatiletech/sqlboiler/v4
  go get github.com/volatiletech/null/v8
  GO111MODULE=off go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
  mv "${GOPATH}/bin/sqlboiler-psql" .
}

push_back() {
  docker push gcr.io/talkiewalkie-305117/talkiewalkie-back:latest
}

kube_back() {
  kubectl create deployment talkiewalkie-back --image=gcr.io/talkiewalkie-305117/talkiewalkie-back:latest
}

install_proto_plugins() {
  go get google.golang.org/protobuf/cmd/protoc-gen-go  google.golang.org/grpc/cmd/protoc-gen-go-grpc
}

grpc() {
  # protoc -I=/usr/local/include/google/protobuf -I=protos/ \
  protoc -I=protos/ \
    --go_out=pb --go-grpc_out=pb \
    --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
    protos/audio_proc.proto
}

"$@"
