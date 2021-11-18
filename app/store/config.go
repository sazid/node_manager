package store

import (
	"errors"
	"io"
	"node_manager/app"

	"github.com/pelletier/go-toml"
)

var (
	ErrMinGreaterThanMax = errors.New("minimum is greater than maximum")
	ErrNegativeInt       = errors.New("numbers cannot be negative")
	ErrFailedToLoad      = errors.New("failed to read config data")
)

type Config struct {
	server           string
	apiKey           string
	minNodes         int
	maxNodes         int
	nodeIdPrefix     string
	externalServices []app.Service
}

func New() Config {
	return Config{
		server:           "",
		apiKey:           "",
		minNodes:         1,
		maxNodes:         1,
		nodeIdPrefix:     "node",
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
	serverInfo := tree.Get("server").(*toml.Tree)
	url := serverInfo.Get("url").(string)
	apiKey := serverInfo.Get("api").(string)

	c.server = url
	c.apiKey = apiKey

	nodeInfo := tree.Get("nodes").(*toml.Tree)
	minNodes := int(nodeInfo.Get("minimum").(int64))
	maxNodes := int(nodeInfo.Get("maximum").(int64))
	nodeIdPrefix := nodeInfo.Get("node_id_prefix")

	if minNodes < 0 || maxNodes < 0 {
		return ErrNegativeInt
	}
	if minNodes > maxNodes {
		return ErrMinGreaterThanMax
	}
	if nodeIdPrefix == nil || nodeIdPrefix.(string) == "" {
		nodeIdPrefix = "nmg"
	}

	c.minNodes = minNodes
	c.maxNodes = maxNodes
	c.nodeIdPrefix = nodeIdPrefix.(string)

	return nil
}

func (c *Config) MaxNodes() int {
	return c.maxNodes
}

func (c *Config) MinNodes() int {
	return c.minNodes
}

func (c *Config) Server() string {
	return c.server
}

func (c *Config) APIKey() string {
	return c.apiKey
}

func (c *Config) NodeIdPrefix() string {
	return c.nodeIdPrefix
}

//func (c *Config) ExternalServices() []app.Service {
//	return c.externalServices
//}
