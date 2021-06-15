package poll_node_state

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"log"
	"path/filepath"
)

const (
	statusFileName = "status"

	statusComplete   = "complete"
	statusIdle       = "idle"
	statusInProgress = "in_progress"
)

type Result struct {
	Complete   int
	Idle       int
	InProgress int
}

// Service polls the file system to see which nodes are currently `Idle`
// or `InProgress`. There is a third intermediate state - `Completed`
// which indicates that a node has completed execution of a test run)
// and now awaiting the next command. Another service can then look for
// services which are in the completed state and upload reports for them.
// State flow:
//
// Start node -> `Idle` -> `InProgress` -> `Completed` -> Kill node
type Service struct {
	fsys fs.FS
}

func (s *Service) Run(context.Context, interface{}) (result interface{}, err error) {
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
		dirEntries, err := fs.ReadDir(s.fsys, nodeDir.Name())
		if err != nil {
			log.Println("err: failed to open node directory.", err)
			continue
		}

		if !fileExistsInDir(dirEntries, statusFileName) {
			continue
		}

		statusFile, err := s.fsys.Open(filepath.Join(
			nodeDir.Name(), statusFileName))
		if err != nil {
			return Result{}, err
		}

		status := readStatus(statusFile)

		switch status {
		case statusInProgress:
			inProgress++
		case statusIdle:
			idle++
		case statusComplete:
			complete++
		}

		_ = statusFile.Close()
	}

	return Result{
		Complete:   complete,
		Idle:       idle,
		InProgress: inProgress,
	}, nil
}

func fileExistsInDir(dirEntries []fs.DirEntry, fileName string) bool {
	for _, f := range dirEntries {
		if f.Name() == fileName {
			return true
		}
	}
	return false
}

func readStatus(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text()
}
