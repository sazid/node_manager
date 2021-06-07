package app

import (
	"context"
	"errors"
	"testing"
	"time"
)

type spyTimerFunc struct {
	Interval time.Duration
	Services []Service
}

func (s *spyTimerFunc) Run(ctx context.Context, _ interface{}) error {
	ticker := time.NewTicker(s.Interval)
	for {
		select {
		case <-ctx.Done():
			return errors.New("timer cancelled")
		case <-ticker.C:
			for _, srv := range s.Services {
				//goland:noinspection GoUnhandledErrorResult
				go srv.Run(ctx)
			}
		}
	}
}

func TestServiceTimer(t *testing.T) {
	t.Run("start and run a timer", func(t *testing.T) {
		var timer Timer = NewServiceTimer(1 * time.Millisecond)

		err := timer.Run(context.Background(), nil)

		if err != nil {
			t.Errorf("got error, did not expect one")
		}
	})

	t.Run("start and stop a timer", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var timer Timer = NewServiceTimer(100 * time.Millisecond)
		ch := make(chan struct{}, 1)

		go func() {
			err := timer.Run(ctx, nil)
			if err == nil {
				t.Errorf("expected an error, got none")
			}
			ch <- struct{}{}
		}()
		time.Sleep(10 * time.Millisecond)
		cancel()
		<-ch
	})

	t.Run("run a service at least 3 times", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch := make(chan struct{}, 1)

		calls := 0
		const minCalls = 3
		var srv ServiceFunc = func(ctx context.Context) error {
			calls++
			return nil
		}

		var timer Timer = &spyTimerFunc{
			Interval: 2 * time.Millisecond,
			Services: []Service{srv},
		}

		go func() {
			_ = timer.Run(ctx, nil)
			if calls < minCalls {
				t.Errorf("got %d calls, want at minimum %d calls", calls, minCalls)
			}

			ch <- struct{}{}
		}()

		time.Sleep(10 * time.Millisecond)
		cancel()
		<-ch
	})

}
