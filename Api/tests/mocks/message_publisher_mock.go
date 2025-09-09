package mocks

type MockMessagePublisher struct {
	PublishFunc func(queue string, body []byte) error
	CloseFunc   func() error
	Messages    []PublishedMessage
}

type PublishedMessage struct {
	Queue string
	Body  []byte
}

func (m *MockMessagePublisher) Publish(queue string, body []byte) error {
	m.Messages = append(m.Messages, PublishedMessage{Queue: queue, Body: body})
	if m.PublishFunc != nil {
		return m.PublishFunc(queue, body)
	}
	return nil
}

func (m *MockMessagePublisher) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}