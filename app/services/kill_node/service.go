package kill_node

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"log"
	"node_manager/app"
	"node_manager/app/services/node_remover"
	"os"
	"path/filepath"
	"strconv"
)

const (
	pidFileName = "pid.txt"

	// PIDs are non-negative in either unix (linux/mac) or windows systems.
	// See https://stackoverflow.com/a/10019054 (unix)
	// and https://stackoverflow.com/a/46058651 (windows)
	pidSentinelValue = -99999999
)

type Service struct {
	fsys        fs.FS
	nodeRemover app.Service
}

func New(fsys fs.FS, nodeRemover app.Service) Service {
	return Service{
		fsys:        fsys,
		nodeRemover: nodeRemover,
	}
}

func (s Service) Run(ctx context.Context, _ interface{}) (result interface{}, err error) {
	// Steps to kill a node
	// 1. Iterate over all the node dirs and see which nodes are in `idle` state.
	// 2. Get the `pid.txt` file and kill the node with the given PID.
	// 3. Call `remove node` service to remove the node from the file system.

	nodesDir, err := fs.ReadDir(s.fsys, ".")
	if err != nil {
		return nil, err
	}

	for _, nodeDir := range nodesDir {
		dirEntries, err := fs.ReadDir(s.fsys, nodeDir.Name())
		if err != nil {
			log.Printf("warn: failed to open node directory: `%s`, continuing.\n%s", nodeDir.Name(), err)
			continue
		}

		if !app.FileExistsInDir(dirEntries, pidFileName) {
			log.Printf("info: `%s` does not exist in the node at `%s`", pidFileName, nodeDir.Name())
			continue
		}

		pidFile, err := s.fsys.Open(filepath.Join(
			nodeDir.Name(), pidFileName))
		if err != nil {
			log.Printf("warn: failed to open %s file for reading node pid.", pidFileName)
			continue
		}

		pid, err := readNodePID(pidFile)
		if err != nil {
			log.Println("warn: failed to read node pid.", err)
			continue
		}

		if err := killProcess(pid); err != nil {
			return nil, err
		}

		msg := node_remover.Message{
			NodeAbsolutePath: filepath.Join(nodeDir.Name(), pidFileName),
		}
		if _, err = s.nodeRemover.Run(ctx, msg); err != nil {
			return nil, err
		}
	}

	return
}

// readNodePID reads the PID from the `io.Reader` and returns
// it as an `int`. The value will be set to `pidSentinelValue`
// in the event of an error.
func readNodePID(r io.Reader) (pid int, err error) {
	scan := bufio.NewScanner(r)
	scan.Scan()
	pid, err = strconv.Atoi(scan.Text())
	return
}

// killProcess finds the process with the given PID and kills it.
// This ignores pid with value `pidSentinelValue`.
func killProcess(pid int) error {
	if pid == pidSentinelValue {
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err = proc.Kill(); err != nil {
		return err
	}
	return nil
}
