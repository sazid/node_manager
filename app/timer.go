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

func NewServiceTimer(interval time.Duration) *ServiceTimer {
	return &ServiceTimer{
		Interval: interval,
		Services: make([]Service, 0),
	}
}

func (s *ServiceTimer) Run(ctx context.Context, _ interface{}) error {
	ticker := time.NewTicker(s.Interval)
	for {
		select {
		case <-ctx.Done():
			return ErrTimerCancelled
		case <-ticker.C:
			return nil
		}
	}
}
