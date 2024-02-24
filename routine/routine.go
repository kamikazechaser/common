package routine

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

const defaultGracefulShutdownPeriod = time.Second * 15

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
		logg: o.Logg,
	}

	if o.GracefulShutdownPeriod == 0 {
		routineManager.gracefulShutdownPeriod = defaultGracefulShutdownPeriod
	}

	return routineManager
}

func (m *RoutineManager) RegisterRoutine(routine Routine) {
	m.routines = append(m.routines, routine)
	m.logg.Debug("registered routine", "routine", routine.Name())
}

func (m *RoutineManager) RunAll(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for _, r := range m.routines {
		wg.Add(1)
		go func(r Routine) {
			defer wg.Done()
			m.startRoutine(ctx, r)
		}(r)
		m.logg.Debug("started background routine sucessfully", "routine", r.Name())
	}

	<-ctx.Done()
	if ctx.Err() != nil {
		m.logg.Debug("graceful shutdown triggered")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), m.gracefulShutdownPeriod)
	defer cancel()

	for _, r := range m.routines {
		wg.Add(1)
		go func(r Routine) {
			defer wg.Done()
			m.stopRoutine(shutdownCtx, r)
		}(r)
		m.logg.Debug("background routine shutdown successfully", "routine", r.Name())
	}

	wg.Wait()
	return nil
}

func (m *RoutineManager) startRoutine(ctx context.Context, routine Routine) {
	if err := routine.Start(ctx); err != nil {
		m.logg.Error("error starting routine", "routine", routine.Name(), "error", err)
	}
}

func (m *RoutineManager) stopRoutine(ctx context.Context, routine Routine) {
	if err := routine.Shutdown(ctx); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			m.logg.Error("graceful period exceeded", "routine", routine.Name())
		}
		if !errors.Is(err, context.Canceled) {
			m.logg.Error("error shutting down routine", "routine", routine.Name(), "error", err)
		}
	}
}
