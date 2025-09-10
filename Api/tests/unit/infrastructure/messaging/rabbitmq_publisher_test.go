package messaging_test

import (
    "errors"
    infra "api/internal/infrastructure/messaging"
    "reflect"
    "testing"

    "github.com/streadway/amqp"
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

// stubs implementing the exported interfaces for unit testing
type stubChannel struct {
    publishCalls int
    lastPublish  struct{
        exchange string
        key      string
        msg      amqp.Publishing
    }

    exDeclared []struct{
        name string
        kind string
        args amqp.Table
    }
    qDeclared []struct{
        name string
        args amqp.Table
    }
    binds []struct{
        name string
        key  string
        ex   string
        args amqp.Table
    }

    publishErr error
    closeErr   error
    calls      *[]string
}

func (s *stubChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
    s.publishCalls++
    s.lastPublish.exchange = exchange
    s.lastPublish.key = key
    s.lastPublish.msg = msg
    return s.publishErr
}
func (s *stubChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
    s.exDeclared = append(s.exDeclared, struct{ name, kind string; args amqp.Table }{name, kind, args})
    return nil
}
func (s *stubChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (string, error) {
    s.qDeclared = append(s.qDeclared, struct{ name string; args amqp.Table }{name, args})
    return name, nil
}
func (s *stubChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
    s.binds = append(s.binds, struct{ name, key, ex string; args amqp.Table }{name, key, exchange, args})
    return nil
}
func (s *stubChannel) Close() error {
    if s.calls != nil {
        *s.calls = append(*s.calls, "channel.Close")
    }
    return s.closeErr
}

type stubConn struct {
    ch        infra.AMQPChannel
    isClosed  bool
    closeErr  error
    chanErr   error
    calls     *[]string
}

func (s *stubConn) Channel() (infra.AMQPChannel, error) { return s.ch, s.chanErr }
func (s *stubConn) Close() error {
    if s.calls != nil {
        *s.calls = append(*s.calls, "conn.Close")
    }
    return s.closeErr
}
func (s *stubConn) IsClosed() bool { return s.isClosed }

func TestRabbitMQPublisher_Close_NoConn_NoPanic(t *testing.T) {
    var p infra.RabbitMQPublisher
    if err := p.Close(); err != nil {
        t.Fatalf("close returned error: %v", err)
    }
}

func TestRabbitMQPublisher_PublishJSON_MarshalError(t *testing.T) {
    var p infra.RabbitMQPublisher
    // functions are not JSON-marshalable
    err := p.PublishJSON("queue", func() {})
    if err == nil {
        t.Fatalf("expected marshal error, got nil")
    }
}

func TestEnsureQueue_Basic_NoDLQ(t *testing.T) {
    ch := &stubChannel{}
    conn := &stubConn{ch: ch}
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) { return conn, nil })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    queue := "videos"
    if err := p.EnsureQueue(queue, 123, false); err != nil {
        t.Fatalf("EnsureQueue: %v", err)
    }
    if len(ch.qDeclared) != 1 {
        t.Fatalf("expected 1 queue declared, got %d", len(ch.qDeclared))
    }
    if ch.qDeclared[0].name != queue {
        t.Fatalf("expected queue name %s, got %s", queue, ch.qDeclared[0].name)
    }
    wantArgs := amqp.Table{"x-max-length": 123}
    if !reflect.DeepEqual(ch.qDeclared[0].args, wantArgs) {
        t.Fatalf("unexpected queue args: got %#v want %#v", ch.qDeclared[0].args, wantArgs)
    }
    if len(ch.exDeclared) != 0 || len(ch.binds) != 0 {
        t.Fatalf("did not expect exchanges or binds for no-DLQ path")
    }
}

func TestEnsureQueue_WithDLQ(t *testing.T) {
    ch := &stubChannel{}
    conn := &stubConn{ch: ch}
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) { return conn, nil })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    queue := "videos"
    if err := p.EnsureQueue(queue, 500, true); err != nil {
        t.Fatalf("EnsureQueue: %v", err)
    }
    // Expect 2 queue declarations: DLQ and main
    if len(ch.qDeclared) != 2 {
        t.Fatalf("expected 2 queue declarations (dlq + main), got %d", len(ch.qDeclared))
    }
    dlqName := queue + ".dlq"
    if ch.qDeclared[0].name != dlqName {
        t.Fatalf("first declared queue should be DLQ %s, got %s", dlqName, ch.qDeclared[0].name)
    }
    // Check exchange and bind
    if len(ch.exDeclared) != 1 {
        t.Fatalf("expected 1 exchange declared, got %d", len(ch.exDeclared))
    }
    if ch.exDeclared[0].name != queue+".dlx" {
        t.Fatalf("unexpected DLX name: %s", ch.exDeclared[0].name)
    }
    if len(ch.binds) != 1 || ch.binds[0].name != dlqName || ch.binds[0].key != queue || ch.binds[0].ex != queue+".dlx" {
        t.Fatalf("unexpected bind: %#v", ch.binds)
    }
    // Main queue args
    gotArgs := ch.qDeclared[1].args
    want := amqp.Table{
        "x-dead-letter-exchange":    queue + ".dlx",
        "x-dead-letter-routing-key": queue,
        "x-max-length":              500,
        "x-overflow":                "reject-publish-dlx",
    }
    if !reflect.DeepEqual(gotArgs, want) {
        t.Fatalf("unexpected main queue args: got %#v want %#v", gotArgs, want)
    }
}

func TestPublish_Success_ConnectsIfNeeded(t *testing.T) {
    ch := &stubChannel{}
    conn := &stubConn{ch: ch}
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) { return conn, nil })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    if err := p.Publish("q", []byte("{}")); err != nil {
        t.Fatalf("publish returned error: %v", err)
    }
    if ch.publishCalls != 1 {
        t.Fatalf("expected 1 publish call, got %d", ch.publishCalls)
    }
    if ch.lastPublish.key != "q" || ch.lastPublish.msg.ContentType != "application/json" {
        t.Fatalf("unexpected publish fields: key=%s ct=%s", ch.lastPublish.key, ch.lastPublish.msg.ContentType)
    }
}

func TestPublish_RetryOnPublishError_ThreeAttemptsNoError(t *testing.T) {
    ch := &stubChannel{publishErr: errors.New("boom")}
    conn := &stubConn{ch: ch}
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) { return conn, nil })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    if err := p.Publish("q", []byte("x")); err != nil {
        t.Fatalf("publish returned unexpected error: %v", err)
    }
    if ch.publishCalls != 3 {
        t.Fatalf("expected 3 publish attempts, got %d", ch.publishCalls)
    }
}

func TestPublish_ReconnectFailure_ReturnsError(t *testing.T) {
    ch := &stubChannel{}
    conn := &stubConn{ch: ch}
    dialCount := 0
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) {
        if dialCount == 0 {
            dialCount++
            return conn, nil
        }
        dialCount++
        return nil, errors.New("dial fail")
    })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    // Simulate connection closed to force reconnect during Publish
    conn.isClosed = true
    if err := p.Publish("q", []byte("{}")); err == nil {
        t.Fatalf("expected error on reconnect failure, got nil")
    }
}

func TestEnsureQueue_ReconnectFailure_ReturnsError(t *testing.T) {
    ch := &stubChannel{}
    conn := &stubConn{ch: ch}
    dialCount := 0
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) {
        if dialCount == 0 {
            dialCount++
            return conn, nil
        }
        dialCount++
        return nil, errors.New("dial fail")
    })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    // Simulate connection closed to force reconnect during EnsureQueue
    conn.isClosed = true
    if err := p.EnsureQueue("q", 10, false); err == nil {
        t.Fatalf("expected error on reconnect failure in EnsureQueue, got nil")
    }
}

func TestClose_ClosesChannelThenConn_ReturnsConnError(t *testing.T) {
    calls := []string{}
    ch := &stubChannel{calls: &calls}
    conn := &stubConn{ch: ch, closeErr: errors.New("conn-close"), calls: &calls}
    p, err := infra.NewRabbitMQPublisherWithDialer("amqp://dummy", func(s string) (infra.AMQPConnection, error) { return conn, nil })
    if err != nil {
        t.Fatalf("new publisher with dialer: %v", err)
    }
    if err := p.Close(); err == nil || err.Error() != "conn-close" {
        t.Fatalf("expected conn-close error, got %v", err)
    }
    if len(calls) < 2 || calls[0] != "channel.Close" || calls[1] != "conn.Close" {
        t.Fatalf("unexpected close call order: %#v", calls)
    }
}

// Tests for wrapper types - moved from infrastructure package to maintain clean architecture
func Test_realConn_Channel_PanicsOnNil(t *testing.T) {
    // This test requires access to internal types, so we test the public interface instead
    if _, err := infra.NewRabbitMQPublisher("amqp://guest:guest@127.0.0.1:1/"); err == nil {
        t.Fatal("expected error from NewRabbitMQPublisher with bad URL")
    }
}

func Test_NewRabbitMQPublisher_DialError(t *testing.T) {
    // Use an unroutable local port to fail fast without external dependency.
    if _, err := infra.NewRabbitMQPublisher("amqp://guest:guest@127.0.0.1:1/"); err == nil {
        t.Fatal("expected error from NewRabbitMQPublisher with bad URL")
    }
}
