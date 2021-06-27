package node_remover

import (
	"context"
	"os"
)

type Message struct {
	Dir string
}

// Service will remove the specified directory.
type Service struct{}

func New() Service {
	return Service{}
}

func (s *Service) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	msg := message.(Message)

	err = os.RemoveAll(msg.Dir)

	return
}
