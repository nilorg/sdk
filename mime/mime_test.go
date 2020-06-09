package mime

import (
	"testing"
)

func TestLookup(t *testing.T) {
	mimeType, exist := Lookup(".jpg")
	t.Log(mimeType, exist)
	mimeType, exist = Lookup("jpg")
	t.Log(mimeType, exist)
}
