package store

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io"
	"node_manager/app"
)

type Config struct {
	minNodes         int
	maxNodes         int
	externalServices []app.Service
}

func NewConfig() Config {
	return Config{
		minNodes:         1,
		maxNodes:         1,
		externalServices: []app.Service{},
	}
}

func (f *Config) Load(reader io.Reader) error {
	tree, err := toml.LoadReader(reader)
	if err != nil {
		return fmt.Errorf("error ocurred while loading config data\n%v", err)
	}

	nodeInfo := tree.Get("nodes").(*toml.Tree)
	f.minNodes = int(nodeInfo.Get("minimum").(int64))
	f.maxNodes = int(nodeInfo.Get("maximum").(int64))

	return nil
}

func (f *Config) MaxNodes() int {
	return f.maxNodes
}

func (f *Config) MinNodes() int {
	return f.minNodes
}

func (f *Config) ExternalServices() []app.Service {
	return f.externalServices
}
