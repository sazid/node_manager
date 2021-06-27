package kill_node

import (
	"context"
	"io/fs"
	"log"
	"node_manager/app"
	"node_manager/app/services/node_remover"
	"path"
	"path/filepath"
)

type Service struct {
	fsys         fs.FS
	nodesDirPath string
	nodeRemover  app.Service
}

func New(fsys fs.FS, nodesDir string, nodeRemover app.Service) Service {
	return Service{
		fsys:         fsys,
		nodesDirPath: nodesDir,
		nodeRemover:  nodeRemover,
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

		//Find the node state file.
		if !app.FileExistsInDir(dirEntries, app.NodeStateFilename) {
			log.Printf("info: `%s` does not exist in the node at `%s`", app.NodeStateFilename, nodeDir.Name())
			continue
		}

		stateFile, err := s.fsys.Open(path.Join(
			nodeDir.Name(), app.NodeStateFilename))
		if err != nil {
			log.Printf("warn: failed to open %s file for reading node state.", app.NodeStateFilename)
			continue
		}

		state, err := app.ReadNodeState(stateFile)
		if err != nil {
			log.Println("warn: failed to read node state.", err)
			_ = stateFile.Close()
			continue
		}
		_ = stateFile.Close()

		switch state {
		case app.StateIdle:
		// TODO (report): Remove the app.StateComplete case once the report
		// service is added.
		case app.StateComplete:
		default:
			log.Printf("info: skipping node `%s` with state `%s`", nodeDir.Name(), state)
			continue
		}

		//Find the PID file and kill node.
		if !app.FileExistsInDir(dirEntries, app.PidFilename) {
			log.Printf("info: `%s` does not exist in the node at `%s`", app.PidFilename, nodeDir.Name())
			continue
		}

		pidFile, err := s.fsys.Open(path.Join(
			nodeDir.Name(), app.PidFilename))
		if err != nil {
			log.Printf("warn: failed to open %s file for reading node pid.", app.PidFilename)
			continue
		}

		pid, err := app.ReadNodePID(pidFile)
		if err != nil {
			log.Println("warn: failed to read node pid.", err)
			_ = pidFile.Close()
			continue
		}
		_ = pidFile.Close()

		err = app.KillProcess(pid)
		if err != nil {
			log.Printf("failed to kill node with PID: %d, err: %+v", pid, err)
		}

		msg := node_remover.Message{
			Dir: filepath.Join(s.nodesDirPath, nodeDir.Name()),
		}
		if _, err = s.nodeRemover.Run(ctx, msg); err != nil {
			log.Printf("failed to remove node %s", nodeDir.Name())
		}

		// Break out after we kill one node.
		break
	}

	return
}
