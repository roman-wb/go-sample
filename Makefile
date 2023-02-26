include .env

# Development
dev-app-run:
	go run cmd/app/main.go

dev-app-hrl:
	air -c build/dev/air.toml

dev-docker-up:
	docker-compose -f build/dev/docker-compose.yml up -d

dev-docker-down:
	docker-compose -f build/dev/docker-compose.yml down

dev-docker-logs:
	docker-compose -f build/dev/docker-compose.yml logs

# Lint
lint:
	golangci-lint run

# Migrations
migration-create:
	goose -dir migrations postgres $(APP_PG_URL) create $(name) sql

migration-up:
	goose -dir migrations -allow-missing postgres $(APP_PG_URL) up

migration-down:
	goose -dir migrations -allow-missing postgres $(APP_PG_URL) down

migration-redo:
	goose -dir migrations -allow-missing postgres $(APP_PG_URL) redo

migration-status:
	goose -dir migrations -allow-missing postgres $(APP_PG_URL) status

migration-reset:
	goose -dir migrations -allow-missing postgres $(APP_PG_URL) reset

# Test
test:
	go test -count=1 -race -cover -coverprofile=coverage.out ./...

coverage:
	go tool cover -html=coverage.out

# integration-test:
# 	go test --tags=integration -count=1 -race -cover  ./...
