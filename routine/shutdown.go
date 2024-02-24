package routine

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NotifyShutdown() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}
