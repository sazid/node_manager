package kill_node

import (
	"context"
	"io/fs"
	"log"
	"node_manager/app"
	"node_manager/app/services/node_remover"
	"path"
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
			break
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

		if err := app.KillProcess(pid); err != nil {
			return nil, err
		}

		msg := node_remover.Message{
			NodeAbsolutePath: path.Join(nodeDir.Name(), app.PidFilename),
		}
		if _, err = s.nodeRemover.Run(ctx, msg); err != nil {
			return nil, err
		}
	}

	return
}
