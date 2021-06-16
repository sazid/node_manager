package start_node

import (
	"context"
	"node_manager/app/store"
)

// Service starts a node and leaves it running. A node goes into the
// `Idle` state when it first starts and go online.
type Service struct {
	config store.Config
}

func (s *Service) Run(context.Context, interface{}) (result interface{}, err error) {
	return nil, nil
}
