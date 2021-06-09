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
			min int
			max int
			err error
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
			tempFile := DummyConfigFile(t, c.min, c.max)

			err := config.Load(tempFile)
			if err != c.err {
				t.Fatalf("got error %q, want %q", err, c.err)
			}
			if err != nil {
				continue
			}

			if config.MinNodes() != c.min {
				t.Errorf("got minimum nodes %d, want %d", config.MinNodes(), c.min)
			}

			if config.MaxNodes() != c.max {
				t.Errorf("got maximum nodes %d, want %d", config.MaxNodes(), c.max)
			}

			_ = tempFile.Close()
			_ = os.Remove(tempFile.Name())
		}
	})
}

// DummyConfigFile creates a temporary file, writes some dummy data into it
// and then returns the `*os.File`. This file must be closed preferably
// with a `defer os.Remove(file.Name)`.
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
	if err := file.Close(); err != nil {
		t.Fatal("failed to close file", err)
	}

	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatal("failed to open temp config file", err)
	}

	return file
}
