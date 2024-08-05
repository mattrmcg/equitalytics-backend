test:
	@go test -v ./...

test-seed:
	@go test -v ./cmd/data

run:
	@go run cmd/main.go

run-seed:
	@go run cmd/data/main.go cmd/data/unmarshal_types.go cmd/data/update_market.go seed

run-market:
	@go run cmd/data/main.go cmd/data/unmarshal_types.go cmd/data/update_market.go market

run-update:
	@go run cmd/data/main.go cmd/data/unmarshal_types.go cmd/data/update_market.go update

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

.PHONY: build webserver data migrate

build: webserver data migrate

build-server:
	go build -o bin/server ./cmd

build-data:
	go build -o bin/data ./cmd/data

build-migrate:
	go build -o bin/migrate ./cmd/migrate