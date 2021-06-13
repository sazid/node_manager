package poll_node_state

import (
	"context"
	"node_manager/app"
	"node_manager/app/store"
)

// Service decides whether to start any new nodes. If it needs to start
// any new node, it'll run the nodeStarterSrv service.
//
// This service calls another service for checking how many nodes are
// currently active. The other service returns the result via a channel.
type Service struct {
	config         store.Config
	nodeStarterSrv app.Service

	activeNodesSrv app.Service
}

type Message struct {
	activeNodes int
}

func New(config store.Config, nodeStarterSrv, activeNodesSrv app.Service) Service {
	return Service{
		config:         config,
		nodeStarterSrv: nodeStarterSrv,
		activeNodesSrv: activeNodesSrv,
	}
}

func (s *Service) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	m, ok := message.(Message)
	if !ok {
		app.PanicOnInvalidMessage(s, Message{})
	}

	// Logic for determining the number of nodes to spin up:
	// # - number of
	//
	// 1. If, (# active nodes < # minimum nodes)
	// 	  Then, start (# minimum nodes - # active nodes)
	//
	// 2. Else If, (# active nodes + 1 <= # max nodes)
	// 	  Then, start 1 more node

	go s.activeNodesSrv.Run(ctx, nil)

	var newNodesToStart int
	if m.activeNodes < s.config.MinNodes() {
		newNodesToStart = s.config.MinNodes() - m.activeNodes
	} else if m.activeNodes+1 <= s.config.MaxNodes() {
		newNodesToStart = 1
	}

	for i := 0; i < newNodesToStart; i++ {
		_, _ = s.nodeStarterSrv.Run(ctx, nil)
	}

	return
}
