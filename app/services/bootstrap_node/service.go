package bootstrap_node

import (
	"context"
	"node_manager/app/store"
)

type Result struct{}

type Service struct {
	config store.Config
}

func (s *Service) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	return Result{}, nil
}
