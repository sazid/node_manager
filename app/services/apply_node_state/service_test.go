package apply_node_state

import (
	"context"
	"node_manager/app/services/poll_node_state"
	"node_manager/app/store"
	"testing"
)

func TestNodeStartKillRuns(t *testing.T) {
	cases := []struct {
		name string

		minNodes int
		maxNodes int

		idle       int
		inProgress int

		expectedStartCount int
		expectedKillCount  int
	}{
		{
			name:               "should start one node",
			minNodes:           1,
			maxNodes:           2,
			idle:               0,
			inProgress:         0,
			expectedStartCount: 1,
			expectedKillCount:  0,
		},
		{
			name:               "should start two nodes",
			minNodes:           2,
			maxNodes:           2,
			idle:               0,
			inProgress:         0,
			expectedStartCount: 2,
			expectedKillCount:  0,
		},
		// {
		// 	name:               "should start one more node after `minimum` no of nodes is in `InProgress` state.",
		// 	minNodes:           2,
		// 	maxNodes:           5,
		// 	idle:               0,
		// 	inProgress:         2,
		// 	expectedStartCount: 1,
		// 	expectedKillCount:  0,
		// },
		// {
		// 	name:               "should not start any more nodes after max limit reached",
		// 	minNodes:           2,
		// 	maxNodes:           5,
		// 	idle:               0,
		// 	inProgress:         5,
		// 	expectedStartCount: 0,
		// 	expectedKillCount:  0,
		// },
		// {
		// 	name:               "should kill 2 `Idle` nodes",
		// 	minNodes:           2,
		// 	maxNodes:           5,
		// 	idle:               3,
		// 	inProgress:         2,
		// 	expectedStartCount: 0,
		// 	expectedKillCount:  2,
		// },
		// {
		// 	name:               "should kill 1 `Idle` node",
		// 	minNodes:           2,
		// 	maxNodes:           5,
		// 	idle:               2,
		// 	inProgress:         2,
		// 	expectedStartCount: 0,
		// 	expectedKillCount:  1,
		// },
		// {
		// 	name:               "should not kill the only `Idle` node when min no of nodes are `InProgress`",
		// 	minNodes:           2,
		// 	maxNodes:           5,
		// 	idle:               1,
		// 	inProgress:         2,
		// 	expectedStartCount: 0,
		// 	expectedKillCount:  0,
		// },
		{
			name:               "should not start/kill any node when there is at least `minNodes` in `Idle` state",
			minNodes:           2,
			maxNodes:           5,
			idle:               2,
			inProgress:         0,
			expectedStartCount: 0,
			expectedKillCount:  0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := store.DummyConfig(t, c.minNodes, c.maxNodes, "", "")
			spyNodeStarter := new(spyNodeStarterService)
			spyNodeKiller := new(spyNodeKillerService)
			spyActiveNodes := &spyActiveNodesService{
				idle:       c.idle,
				inProgress: c.inProgress,
			}
			srv := New(config, spyNodeStarter, spyNodeKiller, spyActiveNodes)

			if _, err := srv.Run(context.Background(), nil); err != nil {
				t.Fatal("got an error, but did not expect one.", err)
			}

			if spyNodeStarter.called != c.expectedStartCount {
				t.Errorf("got start service called %v times, want %v", spyNodeStarter.called, c.expectedStartCount)
			}

			if spyNodeKiller.called != c.expectedKillCount {
				t.Errorf("got kill service called %v times, want %v", spyNodeKiller.called, c.expectedKillCount)
			}
		})
	}
}

type spyNodeStarterService struct {
	called int
}

func (s *spyNodeStarterService) Run(context.Context, interface{}) (result interface{}, err error) {
	s.called++
	return
}

type spyNodeKillerService struct {
	called int
}

func (s *spyNodeKillerService) Run(context.Context, interface{}) (result interface{}, err error) {
	s.called++
	return
}

type spyActiveNodesService struct {
	idle       int
	inProgress int
}

func (s *spyActiveNodesService) Run(context.Context, interface{}) (result interface{}, err error) {
	return poll_node_state.Result{
		InProgress: s.inProgress,
		Idle:       s.idle,
	}, nil
}
