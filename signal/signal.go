package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// AwaitExit 等待退出
func AwaitExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGHUP)
	fmt.Println("awaiting signal")
	<-sigs
	fmt.Println("exiting")
}
