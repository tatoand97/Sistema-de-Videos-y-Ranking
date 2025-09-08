package testing

import (
	"context"
	"testing"
	"time"
)

// TestContext creates a context with timeout for tests
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// AssertEventually retries assertion until success or timeout
func AssertEventually(t *testing.T, assertion func() bool, timeout time.Duration, interval time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if assertion() {
			return
		}
		time.Sleep(interval)
	}
	t.Fatal("assertion failed within timeout")
}

// MockLogger provides a test logger
type MockLogger struct {
	logs []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{logs: make([]string, 0)}
}

func (m *MockLogger) Info(msg string) {
	m.logs = append(m.logs, "INFO: "+msg)
}

func (m *MockLogger) Error(msg string) {
	m.logs = append(m.logs, "ERROR: "+msg)
}

func (m *MockLogger) GetLogs() []string {
	return m.logs
}

func (m *MockLogger) Clear() {
	m.logs = make([]string, 0)
}