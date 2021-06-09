package store

import (
	"node_manager/app/test_utils"
	"os"
	"testing"
)

func TestFileConfig(t *testing.T) {
	t.Run("load config from file", func(t *testing.T) {
		config := NewConfig()

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
			tempFile := test_utils.DummyConfigFile(t, c.minNodes, c.maxNodes)

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

			_ = tempFile.Close()
			_ = os.Remove(tempFile.Name())
		}
	})
}
