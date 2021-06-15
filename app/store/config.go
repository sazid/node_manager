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

func (c *Config) Load(reader io.Reader) error {
	tree, err := toml.LoadReader(reader)
	if err != nil {
		return ErrFailedToLoad
	}

	if err := c.validateAndLoadNodeInfo(tree); err != nil {
		return err
	}

	// TODO: Add validation and data loading for external services.

	return nil
}

func (c *Config) validateAndLoadNodeInfo(tree *toml.Tree) error {
	nodeInfo := tree.Get("nodes").(*toml.Tree)

	minNodes := int(nodeInfo.Get("minimum").(int64))
	maxNodes := int(nodeInfo.Get("maximum").(int64))

	if minNodes < 0 || maxNodes < 0 {
		return ErrNegativeInt
	}
	if minNodes > maxNodes {
		return ErrMinGreaterThanMax
	}

	c.minNodes = minNodes
	c.maxNodes = maxNodes

	return nil
}

func (c *Config) MaxNodes() int {
	return c.maxNodes
}

func (c *Config) MinNodes() int {
	return c.minNodes
}

//func (c *Config) ExternalServices() []app.Service {
//	return c.externalServices
//}
