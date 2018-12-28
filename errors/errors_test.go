package errors

import "testing"

func TestError(t *testing.T) {
	err := New(100, "测试")
	t.Error(err)
}
func TestRpcError(t *testing.T) {
	err := FormatGRpcError(New(500, "rpc error: code = Unknown desc = 500-测试错误"))
	t.Error(err)
}
