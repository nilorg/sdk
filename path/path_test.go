package path

import (
	"fmt"
	"os"
	"testing"
)

func TestDirs(t *testing.T) {
	var ss []string
	err := Dirs(os.Getenv("GOPATH"), &ss)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", ss)
	t.Logf("%+v\n", ss)
}
