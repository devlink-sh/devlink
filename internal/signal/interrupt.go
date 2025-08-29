package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nReceived interrupt, shutting down...")
}
