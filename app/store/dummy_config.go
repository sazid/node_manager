package store

import (
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"
)

const tomlContent = `
[nodes]
minimum = %d
maximum = %d
`

// DummyConfigFile creates a temporary file, writes some dummy data into it,
// and then returns the `*os.File`.
func DummyConfigFile(t testing.TB, minNodes, maxNodes int) fs.File {
	t.Helper()

	data := fmt.Sprintf(tomlContent, minNodes, maxNodes)
	mapFS := fstest.MapFS{
		"config.toml": {Data: []byte(data)},
	}

	file, _ := mapFS.Open("config.toml")

	return file
}

func HelperLoadDummyConfig(t testing.TB, minNodes, maxNodes int) Config {
	t.Helper()

	file := DummyConfigFile(t, minNodes, maxNodes)
	config := NewConfig()

	_ = config.Load(file)
	_ = file.Close()

	return config
}
