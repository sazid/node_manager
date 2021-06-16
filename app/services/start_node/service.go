package start_node

import (
	"context"
	"log"
	"node_manager/app"
	"node_manager/app/services/bootstrap_node"
	"node_manager/app/store"
	"os/exec"
	"path/filepath"
)

const nodeCliFile = "node_cli.py"

// Service starts a node and leaves it running. A node goes into the
// `Idle` state when it first starts and go online.
type Service struct {
	config           store.Config
	bootstrapNodeSrv app.Service
}

type Result struct {
	Path string
}

func (s *Service) Run(ctx context.Context, _ interface{}) (result interface{}, err error) {
	r, err := s.bootstrapNodeSrv.Run(ctx, nil)
	result = Result{
		Path: r.(bootstrap_node.Result).Path,
	}
	log.Println(result.(Result).Path)
	if err != nil {
		return
	}

	nodeCliPath := filepath.Join(result.(Result).Path, nodeCliFile)

	var nodeCliArgs = []string{
		nodeCliPath,
		"-s",
		"https://github.com",
		"-k",
		"123ABC-456DEF-789GHI-101JKL",
		"--once",
	}

	cmd := exec.CommandContext(ctx, "python3", nodeCliArgs...)
	err = cmd.Start()
	if err != nil {
		return
	}

	return
}
