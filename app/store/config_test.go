package store

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileConfig(t *testing.T) {
	t.Run("load config from file", func(t *testing.T) {
		config := NewConfig()
		tempFile := DummyConfigFile(t)
		//goland:noinspection GoUnhandledErrorResult
		defer os.Remove(tempFile.Name())
		//goland:noinspection GoUnhandledErrorResult
		defer tempFile.Close()

		if err := config.Load(tempFile); err != nil {
			t.Fatal("failed to load config data", err)
		}

		if config.MinNodes() != 2 {
			t.Errorf("got minimum nodes %d, want %d", config.MinNodes(), 2)
		}

		if config.MaxNodes() != 5 {
			t.Errorf("got maximum nodes %d, want %d", config.MaxNodes(), 5)
		}
	})
}

const tomlContent = `
[nodes]
minimum = 2
maximum = 5
`

// DummyConfigFile creates a temporary file, writes some dummy data into it
// and then returns the `*os.File`. This file must be closed preferably
// with a `defer os.Remove(file.Name)`.
func DummyConfigFile(t testing.TB) *os.File {
	t.Helper()
	file, err := ioutil.TempFile("", "temp_config_*")
	if err != nil {
		t.Fatal("failed to create temp config file", err)
	}

	if _, err = file.WriteString(tomlContent); err != nil {
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
