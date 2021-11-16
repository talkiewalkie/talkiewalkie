# Server

## Getting started

- Needs `go1.17` installed, `postgresql` version 11 or above.
- We've forked codegen tools for our purposes: [`sqlboiler`](https://github.com/theo-m/sqlboiler), our ORM,
  and [`genny`](https://github.com/theo-m/genny) for generic slice utils. Install instructions available are on each
  repo.
- Run the following and fill with your info:

```bash
cp .env.sample .env.dev
cp sqlboiler.sample.toml sqlboiler.toml
```

- You're good to go!

```bash
createdb talkiewalkie
go run .
```

For development one can use [`air`](https://github.com/cosmtrek/air) to reload the server on file changes. Once
installed just run `air`.

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
