package app

import (
	"context"
	"fmt"
)

type Service interface {
	Run(ctx context.Context, message interface{}) (result interface{}, err error)
}

// ServiceFunc converts any function that takes a `context.Context` into a `app.Service`.
type ServiceFunc func(ctx context.Context, message interface{}) (result interface{}, err error)

func (s ServiceFunc) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	return s(ctx, message)
}

func PanicOnInvalidMessage(srv, msg interface{}) {
	panic(fmt.Sprintf("in %T, message must have type: %T", srv, msg))
}
