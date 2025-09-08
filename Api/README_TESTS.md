# API Unit Tests

This document describes the unit testing architecture for the API following the same pattern as the workers.

## Test Structure

```
Api/
├── tests/
│   ├── mocks/                    # Mock implementations
│   │   ├── video_repository_mock.go
│   │   ├── video_storage_mock.go
│   │   ├── message_publisher_mock.go
│   │   └── user_repository_mock.go
│   ├── unit/                     # Unit tests
│   │   ├── application/          # Use case tests
│   │   │   ├── auth_use_case_test.go
│   │   │   ├── uploads_use_case_simple_test.go
│   │   │   └── validations_test.go
│   │   └── domain/               # Domain entity tests
│   │       └── video_entity_test.go
│   └── integration/              # Integration tests
│       └── handlers/
│           └── uploads_handler_test.go
```

## Test Coverage

### Domain Layer
- ✅ Video entity validation
- ✅ Video status constants
- ✅ Entity creation and properties

### Application Layer
- ✅ Auth service (login, logout, token validation)
- ✅ Upload use case (list videos, get video by ID)
- ✅ Validation functions (MP4 validation)

### Infrastructure Layer
- ✅ Mock implementations for all interfaces
- ✅ Repository pattern mocks
- ✅ Storage service mocks
- ✅ Message publisher mocks

## Running Tests

### All API tests
```bash
cd Api
go test -v ./tests/unit/...
```

### With coverage
```bash
cd Api
make test-coverage
```

### From project root (includes all workers)
```bash
make test-workers
```

## Test Patterns

Following the same architecture as workers:
1. **Mocks in separate package** - Clean separation of test doubles
2. **Unit tests by layer** - Domain, application, infrastructure
3. **Table-driven tests** - Comprehensive test cases
4. **Dependency injection** - Easy mocking of dependencies
5. **Context-aware testing** - Proper context handling

## Mock Usage

All mocks implement the same interfaces as production code:
- `MockVideoRepository` implements `interfaces.VideoRepository`
- `MockVideoStorage` implements `interfaces.VideoStorage`
- `MockMessagePublisher` implements `interfaces.MessagePublisher`
- `MockUserRepository` implements `interfaces.UserRepository`

## Integration with CI/CD

Tests are integrated into the main Makefile and can be run as part of the worker test suite:
- `make test-workers` - Runs all worker and API tests
- `make coverage-workers` - Generates coverage reports for all components