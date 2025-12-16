package storage

import (
	"context"
	"reflect"
	"testing"
)

func TestMetadata_Get_Set(t *testing.T) {
	md := make(Metadata)
	key := "Content-Type"
	value := "application/json"

	// Test Set
	md.Set(key, value)

	// Test Get
	if got := md.Get(key); got != value {
		t.Errorf("Metadata.Get() = %v, want %v", got, value)
	}

	// Test Get non-existent key
	if got := md.Get("Non-Existent"); got != "" {
		t.Errorf("Metadata.Get() = %v, want empty string", got)
	}
}

func TestIncomingContext(t *testing.T) {
	ctx := context.Background()
	md := Metadata{
		"key1": "value1",
		"key2": "value2",
	}

	// Test NewIncomingContext
	ctxWithMD := NewIncomingContext(ctx, md)

	// Test FromIncomingContext
	gotMD, ok := FromIncomingContext(ctxWithMD)
	if !ok {
		t.Fatal("FromIncomingContext() ok = false, want true")
	}

	if !reflect.DeepEqual(gotMD, md) {
		t.Errorf("FromIncomingContext() = %v, want %v", gotMD, md)
	}

	// Test FromIncomingContext with empty context
	_, ok = FromIncomingContext(ctx)
	if ok {
		t.Error("FromIncomingContext() ok = true, want false")
	}
}
