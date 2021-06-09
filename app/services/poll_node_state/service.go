package poll_node_state

import (
	"context"
	"node_manager/app"
	"node_manager/app/store"
)

// Service decides whether to start any new nodes. If it needs to start
// any new node, it'll run the nodeStarter service service.
type Service struct {
	config      *store.Config
	nodeStarter app.Service
}

func (s *Service) Run(ctx context.Context) (err error) {
	if getCurrentNodes() < s.config.MinNodes() {
		newNodesToStart := s.config.MinNodes() - getCurrentNodes()

		for i := 0; i < newNodesToStart; i++ {
			_ = s.nodeStarter.Run(ctx)
		}
	}
	return
}

// getCurrentNodes returns the number of currently active/online nodes.
func getCurrentNodes() int {
	return 0
}