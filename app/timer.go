package app

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimerCancelled = errors.New("timer cancelled")
)

type Timer interface {
	Run(ctx context.Context, store interface{}) error
}

type ServiceTimer struct {
	Interval time.Duration
	Services []Service
}

func NewServiceTimer(interval time.Duration, services []Service) *ServiceTimer {
	return &ServiceTimer{
		Interval: interval,
		Services: services,
	}
}

func (s *ServiceTimer) Run(ctx context.Context, message interface{}) error {
	ticker := time.NewTicker(s.Interval)
	for {
		select {
		case <-ctx.Done():
			return ErrTimerCancelled
		case <-ticker.C:
			for _, srv := range s.Services {
				//goland:noinspection GoUnhandledErrorResult
				go srv.Run(ctx, message)
			}
		}
	}
}
