package routine

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"
)

const defaultGracefulShutdownPeriod = time.Second * 10

type (
	// Routine describes how background go(routines) should be implemented.
	Routine interface {
		// Name of the (go)routines. Used in log identifiers.
		Name() string
		// Start function of the underlying background service is wrapped in this method.
		Start(context.Context) error
		// IgnoreStartError is an error usually returned by Start when Shutdown is called.
		// E.g. [http.ErrServerClosed]
		IgnoreStartError() error
		// Shutdown function of the underlying background service is wrapped in this method.
		Shutdown(context.Context) error
	}

	// RoutineManagerOpts are configurable options to be passed to the RoutineManager.
	RoutineManagerOpts struct {
		// Logg is a required logger that prints go(routine) status when controlled by the RoutineManager.
		Logg *slog.Logger
		// GracefulShutdownPeriod is the max time to wait before forcefully exiting the entire process after Shutdown is called.
		GracefulShutdownPeriod time.Duration
	}

	RoutineManager struct {
		logg                   *slog.Logger
		gracefulShutdownPeriod time.Duration
		routines               []Routine
	}
)

func NewRoutineManager(o RoutineManagerOpts) *RoutineManager {
	routineManager := &RoutineManager{
		logg:                   o.Logg,
		gracefulShutdownPeriod: defaultGracefulShutdownPeriod,
	}

	if o.GracefulShutdownPeriod != 0 {
		routineManager.gracefulShutdownPeriod = o.GracefulShutdownPeriod
	}

	return routineManager
}

// RegisterRoutine registers a new (go)routine which implements Routine.
func (m *RoutineManager) RegisterRoutine(routine Routine) {
	m.routines = append(m.routines, routine)
	m.logg.Debug("successfully registered routine", "routine", routine.Name())
}

// RunAll will start all routines (order not guaranteed) and waits for context cancellation before triggering a Shutdown.
// Use in conjuction with NotifyShutdown.
func (m *RoutineManager) RunAll(ctx context.Context, stop context.CancelFunc) {
	wg := sync.WaitGroup{}

	for _, r := range m.routines {
		wg.Add(1)
		go func(r Routine) {
			defer wg.Done()
			m.startRoutine(ctx, r)
		}(r)
	}

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), m.gracefulShutdownPeriod)

	for _, r := range m.routines {
		wg.Add(1)
		go func(r Routine) {
			defer wg.Done()
			m.stopRoutine(shutdownCtx, r)
		}(r)
	}

	go func() {
		wg.Wait()
		stop()
		cancel()
		os.Exit(0)
	}()

	<-shutdownCtx.Done()
	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		stop()
		cancel()
		m.logg.Error("graceful shutdown period exceeded, forcefully shutting down")
	}
	os.Exit(1)
}

func (m *RoutineManager) startRoutine(ctx context.Context, routine Routine) {
	if err := routine.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) || !errors.Is(err, routine.IgnoreStartError()) {
			m.logg.Error("error starting routine", "routine", routine.Name(), "error", err)
		}
	}
	m.logg.Debug("successfully started routine", "routine", routine.Name())
}

func (m *RoutineManager) stopRoutine(ctx context.Context, routine Routine) {
	if err := routine.Shutdown(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			m.logg.Error("error shutting down routine", "routine", routine.Name(), "error", err)
		}
	}
	m.logg.Debug("successfully shutdown routine", "routine", routine.Name())
}
