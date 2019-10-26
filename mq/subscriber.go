package mq

import "context"

// SubscribeHandler ...
type SubscribeHandler func(ctx context.Context, msg interface{})

// Subscriber ...
type Subscriber interface {
	Register(topic string, h SubscribeHandler, queue ...string) error
}
