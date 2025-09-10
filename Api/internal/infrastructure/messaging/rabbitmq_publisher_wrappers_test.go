package messaging

import (
    "testing"
)

// helper to assert a function panics
func assertPanics(t *testing.T, f func()) {
    t.Helper()
    defer func() {
        if r := recover(); r == nil {
            t.Fatal("expected panic, got none")
        }
    }()
    f()
}

func Test_realConn_Channel_PanicsOnNil(t *testing.T) {
    var rc realConn // rc.c is nil
    assertPanics(t, func() {
        _, _ = rc.Channel()
    })
}

func Test_realConn_Close_PanicsOnNil(t *testing.T) {
    var rc realConn // rc.c is nil
    assertPanics(t, func() {
        _ = rc.Close()
    })
}

func Test_realConn_IsClosed_PanicsOnNil(t *testing.T) {
    var rc realConn // rc.c is nil
    assertPanics(t, func() {
        _ = rc.IsClosed()
    })
}

func Test_realChannel_QueueDeclare_PanicsOnNil(t *testing.T) {
    var ch realChannel // ch.ch is nil
    assertPanics(t, func() {
        _, _ = ch.QueueDeclare("q", true, false, false, false, nil)
    })
}

// Also cover NewRabbitMQPublisher dial path that fails quickly.
func Test_NewRabbitMQPublisher_DialError(t *testing.T) {
    // Use an unroutable local port to fail fast without external dependency.
    if _, err := NewRabbitMQPublisher("amqp://guest:guest@127.0.0.1:1/"); err == nil {
        t.Fatal("expected error from NewRabbitMQPublisher with bad URL")
    }
}
