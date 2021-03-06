package poll_node_state

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"node_manager/app"
	"path"
)

var (
	ErrServiceCancelled = errors.New("`poll_node_state` service cancelled")
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

func New(fsys fs.FS) Service {
	return Service{fsys}
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

		if !app.FileExistsInDir(dirEntries, app.NodeStateFilename) {
			log.Printf("info: `%s` does not exist in the node at `%s`", app.NodeStateFilename, nodeDir.Name())
			// Nodes which were created very recently may not have the node_state.json file created
			idle++
			continue
		}

		statusFile, err := s.fsys.Open(path.Join(
			nodeDir.Name(), app.NodeStateFilename))
		if err != nil {
			log.Printf("warn: failed to open %s file for reading node state.", app.NodeStateFilename)
			continue
		}

		state, err := app.ReadNodeState(statusFile)
		if err != nil {
			log.Println("warn: failed to read node state.", err)
			_ = statusFile.Close()
			continue
		}
		_ = statusFile.Close()

		switch state {
		case app.StateInProgress:
			inProgress++
		case app.StateIdle:
			idle++
		case app.StateComplete:
			complete++
		case app.StateStarting:
			// We consider a node that is starting just now as an `idle` node
			// because we don't want to keep spawning new nodes until we've
			// reached the max limit. We only want 1 more idle node to be
			// launched.
			idle++
		default:
			log.Printf("err: invalid node state: %s", state)
		}
	}

	return Result{
		Complete:   complete,
		Idle:       idle,
		InProgress: inProgress,
	}, nil
}
