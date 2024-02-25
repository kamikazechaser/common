package routine

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// NotifyShutdown is a named utility for signal.NotifyContext
func NotifyShutdown() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}
