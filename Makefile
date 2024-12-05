.PHONY: all
all: ;

.PHONY: test-integration
test-integration:
	go test ./... -v -count=1 -p=1 -tags="integration_tests"

