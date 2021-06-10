package poll_node_state

import (
	"context"
	"fmt"
	"node_manager/app/store"
	"node_manager/app/test_utils"
	"testing"
)

type spyNodeStarterService struct {
	called int
}

func (s *spyNodeStarterService) Run(context.Context) error {
	s.called++
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
			config := LoadDummyConfig(t, c.minNodes, c.maxNodes)
			spyNodeStarter := &spyNodeStarterService{c.called}
			srv := Service{
				config:      config,
				nodeStarter: spyNodeStarter,
			}

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

func LoadDummyConfig(t testing.TB, minNodes, maxNodes int) *store.Config {
	file := test_utils.DummyConfigFile(t, minNodes, maxNodes)
	config := store.NewConfig()

	if err := config.Load(file); err != nil {
		t.Fatal("failed to load config", err)
	}

	if err := file.Close(); err != nil {
		t.Fatal("failed to close temp config file", err)
	}
	return &config
}
