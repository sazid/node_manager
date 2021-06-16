package bootstrap_node

import (
	"context"
	"fmt"
	copy2 "github.com/otiai10/copy"
	"node_manager/app/constants"
	"node_manager/app/store"
	"os"
	"path/filepath"
)

type Result struct {
	Path string
}

const nodeDirPrefix = "node_"

var skipList = []string{
	".git",
	".gitignore",
	".github",
	".vscode",
	".style.yapf",
	"__pycache__",
	"pid.txt",
	"node_id.conf",
	"AutomationLog",
}

// Service copies an already available node into the "nodes" directory
// and returns the path to that directory.
type Service struct {
	config   store.Config
	nodeDir  string
	copyDest string
}

func (s *Service) Run(_ context.Context, _ interface{}) (result interface{}, err error) {
	// 0. [Optional for now] Download node based on config provided version.
	// 1. Copy node
	// 2. Return path to copied node folder.

	dest := generateDirName(s.copyDest)

	if err := os.MkdirAll(dest, constants.OS_USER_RWX); err != nil {
		return Result{}, err
	}

	if err := copyNode(s.nodeDir, dest); err != nil {
		return Result{}, err
	}

	return Result{
		Path: dest,
	}, nil
}

// generateDirName generates a new random directory name with the
// prefix `nodeDirPrefix` and an integer which starts from 0.
//
// https://gist.github.com/mattes/d13e273314c3b3ade33f
func generateDirName(baseDest string) string {
	dest := baseDest
	_, err := os.Stat(dest)
	for i := 0; !os.IsNotExist(err); i++ {
		dest = filepath.Join(
			baseDest,
			fmt.Sprintf("%s%d", nodeDirPrefix, i))
		_, err = os.Stat(dest)
	}
	return dest
}

func copyNode(src, dest string) error {
	opt := copy2.Options{
		Skip: func(src string) (bool, error) {
			for _, s := range skipList {
				src = filepath.Base(src)
				if src == s {
					return true, nil
				}
			}
			return false, nil
		},
	}
	return copy2.Copy(src, dest, opt)
}
