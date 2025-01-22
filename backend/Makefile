GO := go
GO_TEST_FLAGS := -v
COVERAGE_FILE := coverage.out

all: test

test:
	@echo "Running all tests..."
	$(GO) test $(GO_TEST_FLAGS) ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE)

clean:
	@echo "Cleaning up..."
	rm -f $(COVERAGE_FILE)

help:
	@echo "Available targets:"
	@echo "  all          - Run all tests (default target)"
	@echo "  test         - Run all tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean up generated files"
	@echo "  help         - Display this help message"

.PHONY: all test test-coverage clean help