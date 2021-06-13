package app

import (
	"context"
	"testing"
	"time"
)

func TestServiceTimer(t *testing.T) {
	t.Run("start and run a timer", func(t *testing.T) {
		var timer Timer = NewServiceTimer(1*time.Millisecond, nil)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		//goland:noinspection GoUnhandledErrorResult
		go timer.Run(ctx, nil)
	})

	t.Run("start and stop a timer", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var timer Timer = NewServiceTimer(100*time.Millisecond, nil)
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
		var srv ServiceFunc = func(ctx context.Context, message interface{}) (result interface{}, err error) {
			calls++
			return
		}

		var timer Timer = NewServiceTimer(5*time.Millisecond, []Service{srv})

		go func() {
			_ = timer.Run(ctx, nil)
			if calls < minCalls {
				t.Errorf("got %d calls, want at minimum %d calls", calls, minCalls)
			}

			ch <- struct{}{}
		}()

		// Definitely not the best way to test a timer! Maybe we should create a spy for this?
		time.Sleep(100 * time.Millisecond)
		cancel()
		<-ch
	})

}
