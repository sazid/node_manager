package poll_node_state

import (
	"bufio"
	"context"
	"io/fs"
	"log"
	"path/filepath"
)

type Result struct {
	InProgress int
	Idle       int
}

const (
	nodesDirName   = "nodes"
	statusFileName = "status"

	statusIdle       = "idle"
	statusInProgress = "in_progress"
	statusComplete   = "complete"
)

// Service polls the file system to see which nodes are currently `Idle`
// or `InProgress`. There is a third intermediate state - `Completed`
// which indicates that a node has completed execution of a test run)
// and now awaiting the next command. Another service can then look for
// services which are in the completed state and upload reports for them.
// State flow:
//
// Start node -> `Idle` -> `InProgress` -> `Completed` -> Kill node
type Service struct {
	fs fs.FS
}

func (s *Service) Run(_ context.Context, _ interface{}) (result interface{}, err error) {
	nodesDir, err := fs.ReadDir(s.fs, ".")
	if err != nil {
		return Result{}, err
	}

	inProgress := 0

	for _, node := range nodesDir {
		dir, err := fs.ReadDir(s.fs, node.Name())
		if err != nil {
			log.Println("err: failed to open directory.", err)
			continue
		}

		for _, f := range dir {
			if f.Name() == statusFileName {
				path := filepath.Join(node.Name(), f.Name())
				statusFile, err := s.fs.Open(path)
				if err != nil {
					return Result{}, err
				}

				scanner := bufio.NewScanner(statusFile)

				scanner.Scan()
				status := scanner.Text()

				switch status {
				case statusInProgress:
					inProgress++
				}

				if err := statusFile.Close(); err != nil {
					return Result{}, nil
				}
			}
		}
	}

	return Result{
		InProgress: inProgress,
	}, nil
}
