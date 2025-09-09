.PHONY: test test-unit test-coverage clean

# Run all tests
test:
	go test ./tests/... -v

# Run unit tests only
test-unit:
	go test ./tests/unit/... -v

# Run tests with coverage
test-coverage:
	go test ./tests/... -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html

# Install test dependencies
deps:
	go mod tidy
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/require

# Run tests in CI mode
test-ci:
	go test ./tests/... -v -race -coverprofile=coverage.out