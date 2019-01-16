package base62

import "testing"

func TestEncode(t *testing.T) {
	num := 100
	t.Log(Encode(num))
}

func TestDecode(t *testing.T) {
	str := Decode("C1")
	t.Log(str)
}
