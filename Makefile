test:
	@go test -v ./...

test-seed:
	@go test -v ./cmd/data

seed:
	@go run cmd/data/main.go seed