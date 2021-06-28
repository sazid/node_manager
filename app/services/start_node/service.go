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

	// First, check if `python3` is available, if not, we check for `python`. If
	// that's not available, exit and let the user know.
	pythonPath, err := exec.LookPath("python3")
	if err != nil {
		pythonPath, err = exec.LookPath("python")
		if err != nil {
			log.Fatal("could not find `python` or `python3` in PATH")
		}
	}

	cmd := exec.CommandContext(ctx, pythonPath, nodeCliArgs...)
	cmd.Dir = nodePath
	log.Printf("Starting node:\n%s", cmd)

	// TODO: The output writer should be taken as a message
	// or point to current node's AutomationLog dir.
	var outputPipe io.ReadCloser
	if s.outputWriter != nil {
		outputPipe, err = cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// TODO: The output writer should be taken as a message
	// or point to current node's AutomationLog dir.
	if s.outputWriter != nil {
		_, _ = io.Copy(s.outputWriter, outputPipe)
	}

	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Printf("failed to wait for cmd: %+v", err)
		}
	}()

	result = Result{
		NodePath: nodePath,
	}
	return
}
