package store

import (
	"testing"
)

func TestFileConfig(t *testing.T) {
	t.Run("load config from file", func(t *testing.T) {
		config := New()

		cases := []struct {
			minNodes int
			maxNodes int
			err      error
		}{
			{1, 1, nil},
			{2, 5, nil},
			{5, 8, nil},
			{10, 5, ErrMinGreaterThanMax},
			{5, 0, ErrMinGreaterThanMax},
			{-1, 1, ErrNegativeInt},
			// ErrNegativeInt takes higher priority than `ErrMinGreaterThanMax`
			{1, -1, ErrNegativeInt},
		}

		for _, c := range cases {
			tempFile := DummyConfigFile(t, c.minNodes, c.maxNodes, "", "")

			err := config.Load(tempFile)
			if err != c.err {
				t.Fatalf("got error %q, want %q", err, c.err)
			}
			if err != nil {
				continue
			}

			if config.MinNodes() != c.minNodes {
				t.Errorf("got minimum nodes %d, want %d", config.MinNodes(), c.minNodes)
			}

			if config.MaxNodes() != c.maxNodes {
				t.Errorf("got maximum nodes %d, want %d", config.MaxNodes(), c.maxNodes)
			}
		}
	})
}

func TestBadConfig(t *testing.T) {
	t.Run("if bad config is provided, it should retain the previous config", func(t *testing.T) {
		config := New()
		tempFile := DummyConfigFile(t, 5, 1, "", "")

		err := config.Load(tempFile)
		if err == nil {
			t.Fatal("should've errored out.", err)
		}

		if config.MinNodes() != 1 {
			t.Errorf("got minimum no of nodes %d, want %d", config.MinNodes(), 1)
		}

		if config.MaxNodes() != 1 {
			t.Errorf("got maximum no of nodes %d, want %d", config.MaxNodes(), 1)
		}
	})
}
