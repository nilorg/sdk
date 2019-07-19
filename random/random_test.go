package random

import "testing"

func TestNumber(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("Number: %s", Number(6))
	}
}
func TestAZaz09(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("AZaz09: %s", AZaz09(6))
	}
}
