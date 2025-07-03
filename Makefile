run:
	go run cmd/main.go --config=config/config.yml

goose-path:
	export GOOSE_MIGRATION_DIR=internal/migrations

goose-up: goose-path
	goose postgres "postgres://test:test@localhost:5432/test" up

goose-down: goose-path
	goose postgres "postgres://test:test@localhost:5432/test" down

docker-run:
	docker compose up