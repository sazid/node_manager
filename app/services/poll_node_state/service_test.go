package poll_node_state

import (
	"context"
	"fmt"
	"node_manager/app/store"
	"testing"
)

type spyNodeStarterService struct {
	called int
}

func (s *spyNodeStarterService) Run(context.Context) error {
	s.called++
	return nil
}

type spyActiveNodesService struct {
	active     int
	activeChan chan int
}

func (s *spyActiveNodesService) Run(context.Context) error {
	s.activeChan <- s.active
	return nil
}

func TestMinimumNodeStarterRuns(t *testing.T) {
	cases := []struct {
		minNodes int
		maxNodes int
		called   int
	}{
		{
			minNodes: 0,
			maxNodes: 1,
			called:   0,
		},
		{
			minNodes: 1,
			maxNodes: 1,
			called:   0,
		},
		{
			minNodes: 2,
			maxNodes: 2,
			called:   0,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("it should run 'Node Starter' service at least %d times", c.minNodes), func(t *testing.T) {
			config := store.LoadDummyConfig(t, c.minNodes, c.maxNodes)
			spyNodeStarter := &spyNodeStarterService{c.called}
			spyActiveNodes := &spyActiveNodesService{
				active:     0,
				activeChan: make(chan int),
			}
			srv := New(config, spyNodeStarter, spyActiveNodes, spyActiveNodes.activeChan)

			if err := srv.Run(context.Background()); err != nil {
				t.Fatal("got an error, but did not expect one.", err)
			}

			want := config.MinNodes()
			got := spyNodeStarter.called

			if got != want {
				t.Errorf("got service called %v, want %v", got, want)
			}
		})
	}
}

//func TestOneMoreNode(t *testing.T) {
//	config := store.LoadDummyConfig(t, 1, 2)
//
//}
