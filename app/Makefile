.PHONY: build
build:
	go build -o cmd/migrations/migrations ./cmd/migrations

.PHONY: test-integration
test-integration:
	go test ./... -v -count=1 -p=1 -tags="integration_tests"