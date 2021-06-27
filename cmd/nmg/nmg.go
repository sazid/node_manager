package main

import (
	"context"
	"fmt"
	"io/fs"
	"node_manager/app"
	"node_manager/app/constants"
	"node_manager/app/services/apply_node_state"
	"node_manager/app/services/bootstrap_node"
	"node_manager/app/services/kill_node"
	"node_manager/app/services/node_remover"
	"node_manager/app/services/poll_node_state"
	"node_manager/app/services/start_node"
	"node_manager/app/store"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var (
		workingDir, err = os.Getwd()
		configFile      = filepath.Join(workingDir, "config.toml")
		nodesDirPath    = filepath.Join(workingDir, "nodes")
		nodeDir         = filepath.Join(workingDir, "zeuz_node")
	)
	if err != nil {
		panic("cannot get current working directory")
	}

	config := loadConfig(configFile)
	nodesDir := loadDirFS(nodesDirPath) // should this be read from a config file or env variable?

	pollNodeStateSrv := poll_node_state.New(nodesDir)

	nodeRemoverSrv := node_remover.New()
	killNodeSrv := kill_node.New(nodesDir, nodesDirPath, &nodeRemoverSrv)

	bootstrapNodeSrv := bootstrap_node.New(config, nodeDir, nodesDirPath)
	nodeStarterSrv := start_node.New(config, &bootstrapNodeSrv, nil)

	applyNodeStateSrv := apply_node_state.New(config, &nodeStarterSrv, &killNodeSrv, &pollNodeStateSrv)

	setupNodeServicesTimer(&applyNodeStateSrv)
}

func loadConfig(configFile string) store.Config {
	configStore := store.New()
	file, err := os.Open(configFile)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		fmt.Println(err)
		return configStore
	}
	configStore.Load(file)
	return configStore
}

func loadDirFS(nodesDirPath string) fs.FS {
	// Create the dir if not exists.
	if err := os.MkdirAll(nodesDirPath, constants.OS_USER_RWX); err != nil {
		panic(fmt.Sprintf("failed to create the directory for storing generated nodes: %+v", err))
	}
	return os.DirFS(nodesDirPath)
}

func setupNodeServicesTimer(srv app.Service) {
	nodeServicesTimer := app.NewServiceTimer(3*time.Second, []app.Service{srv})
	if err := nodeServicesTimer.Run(context.Background(), nil); err != nil {
		panic(fmt.Sprintf("failed to start timer.%+v", err))
	}
}
