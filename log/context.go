package log

import "context"

type (
	traceIDKey struct{}
	spanIDKey  struct{}
	userIDKey  struct{}
)

// NewTraceIDContext ...
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// FromTraceIDContext ...
func FromTraceIDContext(ctx context.Context) (traceID string, ok bool) {
	traceID, ok = ctx.Value(traceIDKey{}).(string)
	return
}

// NewSpanIDContext ...
func NewSpanIDContext(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey{}, spanID)
}

// FromSpanIDContext ...
func FromSpanIDContext(ctx context.Context) (spanID string, ok bool) {
	spanID, ok = ctx.Value(spanIDKey{}).(string)
	return
}

// NewUserIDContext ...
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// FromUserIDContext ...
func FromUserIDContext(ctx context.Context) (userID string, ok bool) {
	userID, ok = ctx.Value(userIDKey{}).(string)
	return
}
