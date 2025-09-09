package messaging

import (
    "encoding/json"
    "log"
    "time"

    "github.com/streadway/amqp"
)

// AMQPChannel abstracts the subset of methods used from amqp.Channel.
type AMQPChannel interface {
    Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
    ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
    QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (string, error)
    QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
    Close() error
}

// AMQPConnection abstracts the subset of methods used from amqp.Connection.
type AMQPConnection interface {
    Channel() (AMQPChannel, error)
    Close() error
    IsClosed() bool
}

// amqpDialer is a function that dials and returns an AMQPConnection.
type amqpDialer func(url string) (AMQPConnection, error)

// real implementations wrapping streadway/amqp types
type realConn struct{ c *amqp.Connection }

func (r *realConn) Channel() (AMQPChannel, error) {
    ch, err := r.c.Channel()
    if err != nil {
        return nil, err
    }
    return &realChannel{ch: ch}, nil
}
func (r *realConn) Close() error   { return r.c.Close() }
func (r *realConn) IsClosed() bool { return r.c.IsClosed() }

type realChannel struct{ ch *amqp.Channel }

func (r *realChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
    return r.ch.Publish(exchange, key, mandatory, immediate, msg)
}
func (r *realChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
    return r.ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}
func (r *realChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (string, error) {
    q, err := r.ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
    if err != nil {
        return "", err
    }
    return q.Name, nil
}
func (r *realChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
    return r.ch.QueueBind(name, key, exchange, noWait, args)
}
func (r *realChannel) Close() error { return r.ch.Close() }

// RabbitMQPublisher publishes messages to RabbitMQ queues.
type RabbitMQPublisher struct {
    conn    AMQPConnection
    channel AMQPChannel
    url     string
    dial    amqpDialer
}

// NewRabbitMQPublisher constructs a publisher using the real AMQP dialer.
func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
    p := &RabbitMQPublisher{url: url, dial: func(url string) (AMQPConnection, error) {
        c, err := amqp.Dial(url)
        if err != nil {
            return nil, err
        }
        return &realConn{c: c}, nil
    }}
    if err := p.connect(); err != nil {
        return nil, err
    }
    return p, nil
}

// NewRabbitMQPublisherWithDialer allows injecting a custom dialer (useful for tests).
func NewRabbitMQPublisherWithDialer(url string, d func(string) (AMQPConnection, error)) (*RabbitMQPublisher, error) {
    p := &RabbitMQPublisher{url: url, dial: d}
    if err := p.connect(); err != nil {
        return nil, err
    }
    return p, nil
}

func (p *RabbitMQPublisher) connect() error {
    conn, err := p.dial(p.url)
    if err != nil {
        return err
    }
    ch, err := conn.Channel()
    if err != nil {
        _ = conn.Close()
        return err
    }
    p.conn = conn
    p.channel = ch
    return nil
}

func (p *RabbitMQPublisher) isConnected() bool {
    return p.conn != nil && !p.conn.IsClosed() && p.channel != nil
}

// EnsureQueue declares a durable queue with optional DLX and max length.
// It is safe to call multiple times; server will keep existing settings when compatible.
func (p *RabbitMQPublisher) EnsureQueue(queueName string, maxLen int, withDLQ bool) error {
    if !p.isConnected() {
        if err := p.connect(); err != nil {
            return err
        }
    }
    if withDLQ {
        dlxName := queueName + ".dlx"
        dlqName := queueName + ".dlq"
        if err := p.channel.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil); err != nil {
            return err
        }
        if _, err := p.channel.QueueDeclare(dlqName, true, false, false, false, nil); err != nil {
            return err
        }
        if err := p.channel.QueueBind(dlqName, queueName, dlxName, false, nil); err != nil {
            return err
        }
        args := amqp.Table{
            "x-dead-letter-exchange":    dlxName,
            "x-dead-letter-routing-key": queueName,
            "x-max-length":              maxLen,
            "x-overflow":                "reject-publish-dlx",
        }
        if _, err := p.channel.QueueDeclare(queueName, true, false, false, false, args); err != nil {
            return err
        }
        return nil
    }
    // Basic durable queue without DLX
    _, err := p.channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{"x-max-length": maxLen})
    return err
}

// Publish sends a message with persistent delivery mode and application/json content-type.
func (p *RabbitMQPublisher) Publish(queueName string, body []byte) error {
    retries := 3
    for i := 0; i < retries; i++ {
        if !p.isConnected() {
            if err := p.connect(); err != nil {
                if i == retries-1 {
                    return err
                }
                time.Sleep(time.Duration(i+1) * time.Second)
                continue
            }
        }
        err := p.channel.Publish("", queueName, false, false, amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "application/json",
            Body:         body,
        })
        if err == nil {
            return nil
        }
        if i < retries-1 {
            time.Sleep(time.Duration(i+1) * time.Second)
        }
    }
    return nil // Don't fail the operation if messaging fails
}

// PublishJSON marshals v to JSON and publishes it.
func (p *RabbitMQPublisher) PublishJSON(queueName string, v interface{}) error {
    b, err := json.Marshal(v)
    if err != nil {
        return err
    }
    return p.Publish(queueName, b)
}

func (p *RabbitMQPublisher) Close() error {
    if p.channel != nil {
        if err := p.channel.Close(); err != nil {
            log.Printf("rabbitmq publisher channel close: %v", err)
        }
    }
    if p.conn != nil {
        return p.conn.Close()
    }
    return nil
}
