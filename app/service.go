package app

import "context"

type Service interface {
	Run(ctx context.Context) error
}

// ServiceFunc converts any function that takes a `context.Context` into a `app.Service`.
type ServiceFunc func(context.Context) error

func (s ServiceFunc) Run(ctx context.Context) error {
	return s(ctx)
}
