package terminal

import (
	"fmt"
	"os"
)

// Title 设置终端标题
func Title(format string, a ...interface{}) {
	format = fmt.Sprintf("\033]0;%s\007", format)
	fmt.Fprintf(os.Stdout, format, a...)
}
