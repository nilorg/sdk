package strings

import "testing"

func TestPadLeft(t *testing.T) {
	s1 := PadLeft("abc", 3, "1")
	t.Logf("s1:%s\n", s1)
}
func TestPadRight(t *testing.T) {
	s1 := PadRight("abc", 3, "1")
	t.Logf("s1:%s\n", s1)
}
