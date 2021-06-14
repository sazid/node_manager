package app

import (
	"context"
)

type Service interface {
	Run(ctx context.Context, message interface{}) (result interface{}, err error)
}

// ServiceFunc converts any function that takes a `context.Context` into a `app.Service`.
type ServiceFunc func(ctx context.Context, message interface{}) (result interface{}, err error)

func (s ServiceFunc) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	return s(ctx, message)
}
