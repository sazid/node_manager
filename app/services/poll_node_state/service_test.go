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

func (s *spyNodeStarterService) Run(context.Context, interface{}) (result interface{}, err error) {
	s.called++
	return
}

type spyActiveNodesService struct {
	active int
}

func (s *spyActiveNodesService) Run(_ context.Context, msg interface{}) (result interface{}, err error) {
	m, _ := msg.(Message)
	s.active = m.activeNodes
	return
}

func TestMinimumNodeStarterRuns(t *testing.T) {
	cases := []struct {
		minNodes int
		maxNodes int
	}{
		{
			minNodes: 1,
			maxNodes: 1,
		},
		{
			minNodes: 2,
			maxNodes: 2,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("it should run 'Node Starter' service at least %d times", c.minNodes), func(t *testing.T) {
			config := store.LoadDummyConfig(t, c.minNodes, c.maxNodes)
			spyNodeStarter := new(spyNodeStarterService)
			spyActiveNodes := &spyActiveNodesService{
				active: 0,
			}
			var srv Service = New(config, spyNodeStarter, spyActiveNodes)

			if _, err := srv.Run(context.Background(), Message{}); err != nil {
				t.Fatal("got an error, but did not expect one.", err)
			}

			want := config.MinNodes()
			got := spyNodeStarter.called

			if got != want {
				t.Errorf("got service called %v times, want %v", got, want)
			}
		})
	}
}

func TestOneMoreNode(t *testing.T) {
	config := store.LoadDummyConfig(t, 2, 5)
	spyNodeStarter := new(spyNodeStarterService)
	spyActiveNodes := &spyActiveNodesService{
		active: config.MinNodes(),
	}
	srv := New(config, spyNodeStarter, spyActiveNodes)

	if _, err := srv.Run(context.Background(), Message{config.MinNodes()}); err != nil {
		t.Fatal("got an error, but did not expect one.", err)
	}

	// We want 1 more node to be active by now
	want := 1
	got := spyNodeStarter.called

	if got != want {
		t.Errorf("got service called %v times, want %v", got, want)
	}
}

func TestNoMoreNodesAfterMaxLimit(t *testing.T) {
	config := store.LoadDummyConfig(t, 2, 5)
	spyNodeStarter := new(spyNodeStarterService)
	spyActiveNodes := &spyActiveNodesService{
		active: config.MaxNodes(),
	}
	srv := New(config, spyNodeStarter, spyActiveNodes)

	if _, err := srv.Run(context.Background(), Message{config.MaxNodes()}); err != nil {
		t.Fatal("got an error, but did not expect one.", err)
	}

	// Since the number of active nodes is already the maximum allowed,
	// the service should not start any more new nodes.
	want := 0
	got := spyNodeStarter.called

	if got != want {
		t.Errorf("got service called %v times, want %v", got, want)
	}
}
