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
	config           store.Config
	bootstrapNodeSrv app.Service
	outputWriter     io.Writer
}

type Result struct {
	NodePath string
}

func New(config store.Config, bootstrapNodeSrv app.Service, outputWriter io.Writer) Service {
	return Service{
		config:           config,
		bootstrapNodeSrv: bootstrapNodeSrv,
		outputWriter:     outputWriter,
	}
}

func (s *Service) Run(ctx context.Context, _ interface{}) (result interface{}, err error) {
	bootstrapResult, err := s.bootstrapNodeSrv.Run(ctx, nil)
	if err != nil {
		return
	}
	nodePath := bootstrapResult.(bootstrap_node.Result).Path

	nodeCliPath := filepath.Join(nodePath, nodeCliFile)

	var nodeCliArgs = []string{
		nodeCliPath,
		"-s",
		s.config.Server(),
		"-k",
		s.config.APIKey(),
		"--once",
	}

	cmd := exec.CommandContext(ctx, "python", nodeCliArgs...)
	outputPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	log.Printf("Starting node:\n%s", cmd)

	err = cmd.Start()
	if err != nil {
		return
	}

	_, _ = io.Copy(s.outputWriter, outputPipe)
	err = cmd.Wait()

	result = Result{
		NodePath: nodePath,
	}
	return
}
