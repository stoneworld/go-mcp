package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-mcp/transport"
)

func main() {
	t := transport.NewStdioServerTransport()

	go func() {
		if err := t.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan

	userCtx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	serverCtx, cancel := context.WithCancel(userCtx)
	cancel()

	if err := t.Shutdown(userCtx, serverCtx); err != nil {
		fmt.Println(err)
		return
	}
}
