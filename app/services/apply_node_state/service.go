package apply_node_state

import (
	"context"
	"node_manager/app"
	"node_manager/app/services/poll_node_state"
	"node_manager/app/store"
)

// Service decides whether to start/kill any nodes. If it needs to start
// any new node, it'll run the `nodeStarterSrv` service. If it needs to
// kill any node, it'll run the `nodeKillerSrv` service.
//
// This service calls another service for checking how many nodes are
// currently active.
type Service struct {
	config           store.Config
	nodeStarterSrv   app.Service
	nodeKillerSrv    app.Service
	pollNodeStateSrv app.Service
}

func New(config store.Config, nodeStarterSrv, nodeKillerSrv, pollNodeStateSrv app.Service) Service {
	return Service{
		config:           config,
		nodeStarterSrv:   nodeStarterSrv,
		nodeKillerSrv:    nodeKillerSrv,
		pollNodeStateSrv: pollNodeStateSrv,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *Service) Run(ctx context.Context, _ interface{}) (result interface{}, err error) {
	nodeStateCh := make(chan poll_node_state.Result, 1)
	defer close(nodeStateCh)
	go func() {
		res, _ := s.pollNodeStateSrv.Run(ctx, nil)
		nodeStateCh <- res.(poll_node_state.Result)
	}()
	nodeState := <-nodeStateCh

	var (
		// If number of idle nodes is less than the minimum number of nodes, we'll
		// spawn new nodes.
		startCount = max(0, s.config.MinNodes()-nodeState.Idle)

		// Kill all nodes that are done with their work.
		killCount = nodeState.Complete
	)

	for i := 0; i < startCount; i++ {
		_, _ = s.nodeStarterSrv.Run(ctx, nil)
	}

	for i := 0; i < killCount; i++ {
		_, _ = s.nodeKillerSrv.Run(ctx, nil)
	}

	/*
		// OLD LOGIC where it also takes into consideration how many maximum nodes
		// are allowed to run at a time.
		active := nodeState.Idle + nodeState.InProgress

		if active < s.config.MinNodes() {
			startCount = s.config.MinNodes() - active
		} else if active == nodeState.InProgress && active+1 <= s.config.MaxNodes() {
			startCount = 1
		} else if nodeState.Idle > 1 && active > s.config.MinNodes() {
			killCount = nodeState.Idle - 1
		}

		// TODO (report): Remove this block once the report service is added.
		if nodeState.Complete > 0 {
			killCount = nodeState.Complete
		}

		// For killing nodes.

		for i := 0; i < killCount; i++ {
			_, _ = s.nodeKillerSrv.Run(ctx, nil)
		}
	*/

	return
}
