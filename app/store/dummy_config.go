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

// DummyConfigFile creates a temporary file, writes some dummy data into it,
// and then returns the `*os.File`.
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

func LoadDummyConfig(t testing.TB, minNodes, maxNodes int) Config {
	file := DummyConfigFile(t, minNodes, maxNodes)
	config := NewConfig()

	if err := config.Load(file); err != nil {
		t.Fatal("failed to load config", err)
	}

	if err := file.Close(); err != nil {
		t.Fatal("failed to close temp config file", err)
	}
	return config
}

func CleanupFile(f *os.File) {
	_ = f.Close()
	_ = os.Remove(f.Name())
}