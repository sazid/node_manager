package app

import (
	"context"
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

func (s *stubService) Run(ctx context.Context) {
	s.calls = append(s.calls, "start")
	if s.service != nil {
		go s.service.Run(ctx)
	}

	select {
	case <-ctx.Done():
		s.calls = append(s.calls, "stop")
		return
	case <-time.After(s.interval * time.Millisecond):
		s.calls = append(s.calls, "complete")
	}
}

func TestService(t *testing.T) {
	t.Run("start and run to completion", func(t *testing.T) {
		var srv Service = NewStubService(0, nil)

		srv.Run(context.Background())

		got := srv.(*stubService).calls
		want := []string{"start", "complete"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q calls, wanted %q", got, want)
		}
	})

	t.Run("start and cancel service via Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var srv Service = NewStubService(10*time.Millisecond, nil)

		go srv.Run(ctx)

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

		go srv3.Run(ctx)

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
