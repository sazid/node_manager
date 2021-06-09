package poll_node_state

import (
	"context"
	"node_manager/app/store"
	"node_manager/app/test_utils"
	"testing"
)

type spyNodeStarterService struct {
	called int
}

func (s *spyNodeStarterService) Run(_ context.Context) error {
	s.called++
	return nil
}

func TestService(t *testing.T) {
	t.Run("it should run Node starter service at least once", func(t *testing.T) {
		config := LoadDummyConfig(t, 1, 1)
		spyNodeStarter := &spyNodeStarterService{0}
		srv := Service{
			config:      config,
			nodeStarter: spyNodeStarter,
		}

		_ = srv.Run(context.Background())

		want := config.MinNodes()
		got := spyNodeStarter.called

		if got != want {
			t.Errorf("got service called %v, want %v", got, want)
		}
	})

	t.Run("it should run Node starter service at least twice", func(t *testing.T) {
		config := LoadDummyConfig(t, 2, 2)
		spyNodeStarter := &spyNodeStarterService{0}
		srv := Service{
			config:      config,
			nodeStarter: spyNodeStarter,
		}

		_ = srv.Run(context.Background())

		got := spyNodeStarter.called

		if got != config.MinNodes() {
			t.Errorf("got service called %v, want %v", got, config.MinNodes())
		}
	})
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
