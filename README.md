## Sample golang

## Included

- http server (echo)
- postgresql (pgx)
- redis (go-redis)
- fearute flags (flipt)
- migrations (goose)
- graceful shutdown
- json logs (zerolog)
- /pkg - common package ready for export in other services

## Dev requirements

Docker, docker-compose, [air](https://github.com/cosmtrek/air), [golangci-lint](https://github.com/golangci/golangci-lint), [goose](https://github.com/pressly/goose)

## Get started

```
git clone github.com/roman-wb/go-sample
cp .env.sample .env
make dev-docker-up
make dev-app-hrl
```

## Config

Everything config in `.env` file

Priority config: Environment => .env => default

## Makefile commands

- `make dev-app-run` - Run app
- `make dev-app-hrl` - Run app with hot reload (air)
- `make dev-docker-up` - Up docker-compose with deps
- `make dev-docker-up` - Down docker-compose with deps
- `make lint` - Run linter
- `make migration-create name=create_users` - Create goose migration with name `timestamp+create_users` in ./migrations
- `make migration-up` - Up migration
- `make migration-down` - Down migration (one step)
- `make migration-redo` - Redo migration (one step)
- `make migration-status` - Get status migrations
- `make migration-reset` - Reset all migrations
- `make test` - Run tests
- `make coverage` - Open coverage.out in browser
