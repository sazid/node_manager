package start_node

import (
	"context"
	"io"
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
	Config           store.Config
	BootstrapNodeSrv app.Service
	OutputWriter     io.Writer
}

type Result struct {
	NodePath string
}

func (s *Service) Run(ctx context.Context, _ interface{}) (result interface{}, err error) {
	bootstrapResult, err := s.BootstrapNodeSrv.Run(ctx, nil)
	if err != nil {
		return
	}
	nodePath := bootstrapResult.(bootstrap_node.Result).Path

	nodeCliPath := filepath.Join(nodePath, nodeCliFile)

	var nodeCliArgs = []string{
		nodeCliPath,
		"-s",
		s.Config.Server(),
		"-k",
		s.Config.APIKey(),
		"--once",
	}

	cmd := exec.CommandContext(ctx, "python3", nodeCliArgs...)
	outputPipe, err := cmd.StdoutPipe()
	log.Printf("Starting node:\n%s", cmd)

	err = cmd.Start()
	if err != nil {
		return
	}

	_, _ = io.Copy(s.OutputWriter, outputPipe)
	err = cmd.Wait()

	result = Result{
		NodePath: nodePath,
	}
	return
}
