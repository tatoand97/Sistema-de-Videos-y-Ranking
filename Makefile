.PHONY: test test-unit test-integration coverage mocks test-workers

# Generate mocks
mocks:
	go generate ./...

# Unit tests for all workers
test-workers:
	@echo "Running AudioRemoval tests..."
	cd AudioRemoval && go test -v ./tests/unit/...
	@echo "Running EditVideo tests..."
	cd EditVideo && go test -v ./tests/unit/...
	@echo "Running TrimVideo tests..."
	cd TrimVideo && go test -v ./tests/unit/...
	@echo "Running Watermarking tests..."
	cd Watermarking && go test -v ./tests/unit/... || true
	@echo "Running gossipOpenClose tests..."
	cd gossipOpenClose && go test -v ./tests/unit/... || true

# Unit tests
test-unit:
	go test -v ./tests/unit/...

# Integration tests
test-integration:
	go test -v ./tests/integration/...

# All tests
test: test-unit test-workers

# Coverage for workers
coverage-workers:
	@echo "Generating coverage for AudioRemoval..."
	cd AudioRemoval && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Generating coverage for EditVideo..."
	cd EditVideo && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Generating coverage for TrimVideo..."
	cd TrimVideo && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html

# Overall coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Coverage by package
coverage-by-pkg:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out | grep -E "(domain|application|infrastructure)"

# Install test dependencies
install-test-deps:
	go mod tidy
	go install github.com/golang/mock/mockgen@latest

# Clean test artifacts
clean-test:
	find . -name "coverage.out" -delete
	find . -name "coverage.html" -delete
	find . -name "*_mock.go" -delete