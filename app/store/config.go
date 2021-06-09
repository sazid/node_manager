package store

import (
	"errors"
	"github.com/pelletier/go-toml"
	"io"
	"node_manager/app"
)

var (
	ErrMinGreaterThanMax = errors.New("minimum is greater than maximum")
	ErrNegativeInt       = errors.New("numbers cannot be negative")
	ErrFailedToLoad      = errors.New("failed to read config data")
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
		return ErrFailedToLoad
	}

	nodeInfo := tree.Get("nodes").(*toml.Tree)

	minNodes := int(nodeInfo.Get("minimum").(int64))
	maxNodes := int(nodeInfo.Get("maximum").(int64))
	if minNodes < 0 || maxNodes < 0 {
		return ErrNegativeInt
	}
	if minNodes > maxNodes {
		return ErrMinGreaterThanMax
	}
	f.minNodes = minNodes
	f.maxNodes = maxNodes

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
