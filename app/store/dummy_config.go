package store

import (
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"
)

const tomlContent = `
[server]
url = "%s"
api = "%s"

[nodes]
minimum = %d
maximum = %d
`

// DummyConfigFile creates a temporary file, writes some dummy data into it,
// and then returns the `*os.File`.
func DummyConfigFile(t testing.TB, minNodes, maxNodes int, server, apiKey string) fs.File {
	t.Helper()

	data := fmt.Sprintf(tomlContent, server, apiKey, minNodes, maxNodes)
	mapFS := fstest.MapFS{
		"config.toml": {Data: []byte(data)},
	}

	file, _ := mapFS.Open("config.toml")

	return file
}

func DummyConfig(t testing.TB, minNodes, maxNodes int, server, apiKey string) Config {
	t.Helper()

	file := DummyConfigFile(t, minNodes, maxNodes, server, apiKey)
	config := New()

	_ = config.Load(file)
	_ = file.Close()

	return config
}
