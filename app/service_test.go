package app

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

type stubService struct {
	calls    []string
	interval time.Duration
	service  *stubService
}

func NewStubService(interval time.Duration, service *stubService) *stubService {
	return &stubService{
		calls:    make([]string, 0),
		interval: interval,
		service:  service,
	}
}

func (s *stubService) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	s.calls = append(s.calls, "start")
	if s.service != nil {
		go func() {
			_, _ = s.service.Run(ctx, message)
		}()
	}

	select {
	case <-ctx.Done():
		s.calls = append(s.calls, "stop")
	case <-time.After(s.interval * time.Millisecond):
		s.calls = append(s.calls, "complete")
	}
	return
}

func TestService(t *testing.T) {
	t.Run("start and run to completion", func(t *testing.T) {
		var srv Service = NewStubService(0, nil)

		_, _ = srv.Run(context.Background(), nil)

		got := srv.(*stubService).calls
		want := []string{"start", "complete"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q calls, wanted %q", got, want)
		}
	})

	t.Run("start and cancel service via Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var srv Service = NewStubService(10*time.Millisecond, nil)

		go func() {
			_, _ = srv.Run(ctx, nil)
		}()

		// Simulate a situation where we cancel the service after 1ms delay.
		time.Sleep(1 * time.Millisecond)
		cancel()
		time.Sleep(1 * time.Millisecond)

		got := srv.(*stubService).calls
		want := []string{"start", "stop"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q calls, wanted %q", got, want)
		}
	})

	t.Run("cancel a chain of running services", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var (
			srv1 = NewStubService(0*time.Millisecond, nil)
			srv2 = NewStubService(200*time.Millisecond, srv1)
			srv3 = NewStubService(500*time.Millisecond, srv2)
		)

		go func() {
			_, _ = srv3.Run(ctx, nil)
		}()

		// Simulate a situation where we cancel the service after 100ms delay.
		time.Sleep(100 * time.Millisecond)
		cancel()

		// This sleep ensures that the cancel signal is propagated to all the
		// goroutines.
		time.Sleep(10 * time.Millisecond)

		cases := []struct {
			name  string
			calls []string
			want  []string
		}{
			{"service 1", srv1.calls, []string{"start", "complete"}},
			{"service 2", srv2.calls, []string{"start", "stop"}},
			{"service 3", srv3.calls, []string{"start", "stop"}},
		}

		for _, c := range cases {
			if !reflect.DeepEqual(c.calls, c.want) {
				t.Errorf("%q: got %q calls, wanted %q", c.name, c.calls, c.want)
			}
		}
	})
}

func TestServiceFunc(t *testing.T) {
	t.Run("start and run to completion", func(t *testing.T) {
		var srv ServiceFunc = func(ctx context.Context, message interface{}) (result interface{}, err error) {
			log.Println("running service func")
			return
		}

		if _, err := srv.Run(context.Background(), nil); err != nil {
			t.Errorf("got %v, want %v", err, nil)
		}
	})

	t.Run("cancel a running service", func(t *testing.T) {
		var srv ServiceFunc = func(ctx context.Context, message interface{}) (result interface{}, err error) {
			select {
			case <-ctx.Done():
				return nil, nil
			case <-time.After(100 * time.Millisecond):
				return nil, errors.New("service not cancelled")
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := make(chan error)
		go func() {
			_, e := srv.Run(ctx, nil)
			err <- e
		}()

		time.Sleep(1 * time.Millisecond)
		cancel()

		if <-err != nil {
			t.Errorf("expected the service to be cancelled earlier")
		}
	})
}
