package poll_node_state

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"node_manager/app"
	"path/filepath"
)

var (
	ErrServiceCancelled = errors.New("`poll_node_state` service cancelled")
)

const (
	stateFileName = "node_state.json"
)

type Result struct {
	// No of nodes in `complete` state.
	Complete int
	// No of `idle` nodes.
	Idle int
	// No of `in_progress` nodes.
	InProgress int
}

// Service polls the file system to see which nodes are currently `Idle`
// or `InProgress`. There is a third intermediate state - `Complete`
// which indicates that a node has completed execution of a test run)
// and now awaiting the next command. Another service can then look for
// services which are in the completed state and upload reports for them.
// State flow:
//
// Start node -> `Idle` -> `InProgress` -> `Completed` -> Kill node
type Service struct {
	fsys fs.FS
}

func (s *Service) Run(ctx context.Context, _ interface{}) (result interface{}, err error) {
	nodesDir, err := fs.ReadDir(s.fsys, ".")
	if err != nil {
		return Result{}, err
	}

	var (
		complete   = 0
		idle       = 0
		inProgress = 0
	)

	for _, nodeDir := range nodesDir {
		// Allow cancellation of service.
		select {
		case <-ctx.Done():
			return Result{}, ErrServiceCancelled
		default:
		}

		dirEntries, err := fs.ReadDir(s.fsys, nodeDir.Name())
		if err != nil {
			log.Printf("warn: failed to open node directory: `%s`, continuing.\n%s", nodeDir.Name(), err)
			continue
		}

		if !app.FileExistsInDir(dirEntries, stateFileName) {
			log.Printf("info: `%s` does not exist in the node at `%s`", stateFileName, nodeDir.Name())
			continue
		}

		statusFile, err := s.fsys.Open(filepath.Join(
			nodeDir.Name(), stateFileName))
		if err != nil {
			log.Printf("warn: failed to open %s file for reading node state.", stateFileName)
			continue
		}

		state, err := app.ReadNodeState(statusFile)
		if err != nil {
			log.Println("warn: failed to read node state.", err)
			continue
		}

		switch state {
		case app.StateInProgress:
			inProgress++
		case app.StateIdle:
			idle++
		case app.StateComplete:
			complete++
		default:
			log.Printf("err: invalid node state: %s", state)
		}

		_ = statusFile.Close()
	}

	return Result{
		Complete:   complete,
		Idle:       idle,
		InProgress: inProgress,
	}, nil
}
