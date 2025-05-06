run:
	go run cmd/main.go --config=config/config.yml

goose-path:
	export GOOSE_MIGRATION_DIR=internal/adapters/migrations

goose-up:
	goose postgres "postgres://test:test@localhost:5432/test" up

goose-down:
	goose postgres "postgres://test:test@localhost:5432/test" down

docker-run:
	docker compose up