package rsa

import (
	"fmt"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	pub, prv, err := GenRsaKey(2048)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pub)
	fmt.Println(prv)
}
