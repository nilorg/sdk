package terminal

import (
	"fmt"
	"os"
)

// Titlef 设置终端标题
func Titlef(format string, a ...interface{}) {
	format = fmt.Sprintf("\033]0;%s\007", format)
	fmt.Fprintf(os.Stdout, format, a...)
}

// Title 设置终端标题
func Title(title string) {
	Titlef(title)
}
