#!/usr/bin/env bash

start_postgres() {
  pg_ctl -D /opt/homebrew/var/postgresql@12 start
}

install_golang_migrate_cli() {
  brew install golang-migrate
}

fmt() {
  gofmt -w testutils/ api/ common/ cmd/ clients/ repositories/**/* pkg/
  goimports -local github.com/talkiewalkie -w testutils/ api/ common/ cmd/ clients/ repositories/**/* pkg/
}

migrate_up() {
  source .env.dev &&
    migrate -path migrations -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432/talkiewalkie?sslmode=disable up
  go generate
}

nukedb() {
  dropdb talkiewalkie && createdb talkiewalkie
  migrate_up
}

new_migration() {
  migrate create -ext sql -dir migrations/ -seq -digits 3 "$1"
}

new_secret() {
  openssl rand -hex 32
}

install_codegen_tools() {
  # protobuf codegen
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.27.1
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

  pushd ../.. || exit 1
  mkdir -p codegen_tools/
  pushd codegen_tools || exit 1

  # generics codegen
  git clone git@github.com:theo-m/genny.git
  pushd genny || exit 1
  go install
  popd

  # orm codegen
  git clone git@github.com:theo-m/sqlboiler.git
  pushd sqlboiler || exit 1
  ./install-fork.sh
  popd

  popd
  popd
}

push_back() {
  docker push gcr.io/talkiewalkie-305117/talkiewalkie-back:latest
}

kube_back() {
  kubectl create deployment talkiewalkie-back --image=gcr.io/talkiewalkie-305117/talkiewalkie-back:latest
}

install_proto_plugins() {
  go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc
}

grpc() {
  # protoc -I=/usr/local/include/google/protobuf -I=protos/ \
  protoc -I=protos/ \
    --go_out=pb --go-grpc_out=pb \
    --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
    protos/audio_proc.proto
}

help() {
  declare -F
}

"$@"
