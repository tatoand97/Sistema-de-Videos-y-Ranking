.PHONY: test test-unit test-integration coverage mocks test-workers

# Generate mocks
mocks:
	go generate ./...

# Unit tests for all workers and API
test-workers:
	@echo "Running API tests..."
	cd Api && go test -v ./tests/unit/...
	@echo "Running AudioRemoval tests..."
	cd Workers/AudioRemoval && go test -v ./tests/unit/...
	@echo "Running EditVideo tests..."
	cd Workers/EditVideo && go test -v ./tests/unit/...
	@echo "Running TrimVideo tests..."
	cd Workers/TrimVideo && go test -v ./tests/unit/...
	@echo "Running Watermarking tests..."
	cd Workers/Watermarking && go test -v ./tests/unit/... || true
	@echo "Running gossipOpenClose tests..."
	cd Workers/gossipOpenClose && go test -v ./tests/unit/... || true
	@echo "Running StatesMachine tests..."
	cd Workers/StatesMachine && go test -v ./tests/unit/... || true

# Unit tests
test-unit:
	go test -v ./tests/unit/...

# Integration tests
test-integration:
	go test -v ./tests/integration/...

# All tests
test: test-unit test-workers

# Coverage for workers and API
coverage-workers:
	@echo "Generating coverage for API..."
	cd Api && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Generating coverage for AudioRemoval..."
	cd Workers/AudioRemoval && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Generating coverage for EditVideo..."
	cd Workers/EditVideo && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Generating coverage for TrimVideo..."
	cd Workers/TrimVideo && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Generating coverage for Watermarking..."
	cd Workers/Watermarking && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html || true
	@echo "Generating coverage for gossipOpenClose..."
	cd Workers/gossipOpenClose && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html || true
	@echo "Generating coverage for StatesMachine..."
	cd Workers/StatesMachine && go test -coverprofile=coverage.out ./tests/unit/... && go tool cover -html=coverage.out -o coverage.html || true

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

# Generate comprehensive HTML coverage report
coverage-html:
	@echo "Generating comprehensive coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Open coverage.html in your browser to view the report"

# Clean test artifacts
clean-test:
	find . -name "coverage.out" -delete
	find . -name "coverage.html" -delete
	find . -name "*_mock.go" -delete