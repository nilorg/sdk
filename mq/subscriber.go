package mq

import "context"

// SubscribeHandler ...
type SubscribeHandler func(ctx context.Context, data []byte) error

// Subscriber ...
type Subscriber interface {
	Register(topic string, h SubscribeHandler, queue ...string) error
}

type subscriber struct {
	client Clienter
}

func (s *subscriber) Register(topic string, h SubscribeHandler, queue ...string) error {
	return s.client.Subscribe(topic, h, queue...)
}

// NewSubscriber ...
func NewSubscriber(client Clienter) Subscriber {
	return &subscriber{
		client: client,
	}
}
