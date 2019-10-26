package mq

import "context"

// Publisher ...
type Publisher interface {
	Publish(ctx context.Context, subj string, msg interface{}) error
}
