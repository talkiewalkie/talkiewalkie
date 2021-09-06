# Server

## Getting started

Needs `go1.15` installed, `postgresql` version 11 or above.

```bash
createdb talkiewalkie
go run .
```

For development one can use [`air`](https://github.com/cosmtrek/air) (`go get -u github.com/cosmtrek/air`) to reload the server on file changes.

### Getting data

Running the following will insert 500+ walks and a few users to the db.

```bash
go run cmd/faker/main.go
```

## Building docker image remotely

Needs `kubectl` and access to our GKE cluster. Ask @theo-m for help.

```bash
./build.sh
```
