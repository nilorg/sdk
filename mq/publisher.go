package mq

import "context"

// Publisher ...
type Publisher interface {
	Publish(ctx context.Context, msg interface{}) error
}

type publisher struct {
	topic  string
	client Clienter
}

// Publish ...
func (p *publisher) Publish(ctx context.Context, msg interface{}) error {
	return p.client.Publish(ctx, p.topic, msg)
}

// NewPublisher ...
func NewPublisher(topic string, client Clienter) Publisher {
	return &publisher{
		topic:  topic,
		client: client,
	}
}
