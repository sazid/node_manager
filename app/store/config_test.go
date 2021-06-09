package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const tomlContent = `
[nodes]
minimum = %d
maximum = %d
`

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
			tempFile := DummyConfigFile(t, c.minNodes, c.maxNodes)

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

// DummyConfigFile creates a temporary file, writes some dummy data into it,
// closes and reopens it, and then returns the `*os.File`.
func DummyConfigFile(t testing.TB, minNodes, maxNodes int) *os.File {
	t.Helper()
	file, err := ioutil.TempFile("", "temp_config_*")
	if err != nil {
		t.Fatal("failed to create temp config file", err)
	}

	if _, err = file.WriteString(fmt.Sprintf(
		tomlContent,
		minNodes,
		maxNodes,
	)); err != nil {
		t.Fatal("failed to write temp config data to file", err)
	}

	if err := file.Sync(); err != nil {
		t.Fatal("failed to sync/flush content to disk", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal("failed to seek file to origin (0, 0) position", err)
	}

	return file
}
