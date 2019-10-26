package mq

import "context"

// Publisher ...
type Publisher interface {
	Publish(ctx context.Context, msg interface{}) error
}
