package poll_node_state

import (
	"context"
	"errors"
	"node_manager/app"
	"node_manager/app/store"
	"time"
)

var (
	ErrActiveNodesChanTimeout = errors.New("timed out while waiting for active nodes channel")
)

const (
	ActiveNodesChanTimeout = 5 * time.Second
)

// Service decides whether to start any new nodes. If it needs to start
// any new node, it'll run the nodeStarterSrv service.
//
// This service calls another service for checking how many nodes are
// currently active. The other service returns the result via a channel.
type Service struct {
	config         store.Config
	nodeStarterSrv app.Service

	activeNodesSrv  app.Service
	activeNodesChan chan int
}

func New(config store.Config, nodeStarter, activeNodes app.Service, activeNodeChan chan int) Service {
	return Service{
		config:          config,
		nodeStarterSrv:  nodeStarter,
		activeNodesSrv:  activeNodes,
		activeNodesChan: activeNodeChan,
	}
}

func (s *Service) Run(ctx context.Context, message interface{}) (result interface{}, err error) {
	// Logic for determining the number of nodes to spin up:
	// # - number of
	//
	// 1. If, (# active nodes < # minimum nodes)
	// 	  Then, start (# minimum nodes - # active nodes)
	//
	// 2. Else If, (# active nodes + 1 <= # max nodes)
	// 	  Then, start 1 more node

	go s.activeNodesSrv.Run(ctx, nil)

	activeNodes := 0
	select {
	case <-time.After(ActiveNodesChanTimeout):
		return nil, ErrActiveNodesChanTimeout
	case activeNodes = <-s.activeNodesChan:
		break
	}

	var newNodesToStart int
	if activeNodes < s.config.MinNodes() {
		newNodesToStart = s.config.MinNodes() - activeNodes
	} else if activeNodes+1 <= s.config.MaxNodes() {
		newNodesToStart = 1
	}

	for i := 0; i < newNodesToStart; i++ {
		_, _ = s.nodeStarterSrv.Run(ctx, nil)
	}

	return
}
