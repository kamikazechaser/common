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
	Routine interface {
		Name() string
		Start(context.Context) error
		Shutdown(context.Context) error
	}

	RoutineManagerOpts struct {
		Logg                   *slog.Logger
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

func (m *RoutineManager) RegisterRoutine(routine Routine) {
	m.routines = append(m.routines, routine)
}

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
		if !errors.Is(err, context.Canceled) {
			m.logg.Error("error starting routine", "routine", routine.Name(), "error", err)
		}
	}
}

func (m *RoutineManager) stopRoutine(ctx context.Context, routine Routine) {
	if err := routine.Shutdown(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			m.logg.Error("error shutting down routine", "routine", routine.Name(), "error", err)
		}
	}
}
