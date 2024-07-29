package toolbox

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForInterrupt() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
