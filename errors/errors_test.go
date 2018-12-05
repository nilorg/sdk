package errors

import "testing"

func TestError(t *testing.T) {
	err := New(100, "测试")
	t.Error(err)
}
