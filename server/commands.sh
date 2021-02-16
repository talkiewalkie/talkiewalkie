#!/usr/bin/env bash

start_postgres() {
  pg_ctl -D /opt/homebrew/var/postgresql@12 start
}

install_golang_migrate_cli() {
  brew install golang-migrate
}

nukedb() {
  dropdb talkiewalkie && createdb talkiewalkie
  migrate -path migrations -database postgres://theo@localhost:5432/talkiewalkie?sslmode=disable up
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

"$@"
