package bootstrap_node

import (
	"context"
	"node_manager/app/store"
	"node_manager/app/test_utils"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestBootstrapNode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := store.New()

	nodeDir := test_utils.GenerateTempNode(t, skipList)
	defer os.RemoveAll(nodeDir)

	copyDest := test_utils.GenerateTempDir(t, "node_instances")
	defer os.RemoveAll(copyDest)

	srv := Service{
		config:   config,
		nodeDir:  nodeDir,
		copyDest: copyDest,
	}

	nodePath, err := srv.Run(ctx, nil)
	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}

	want := Result{
		Path: filepath.Join(copyDest, "node_0"),
	}
	if !reflect.DeepEqual(nodePath, want) {
		t.Errorf("got %+v, want %+v", nodePath, want)
	}
}

func TestGenerateDirName(t *testing.T) {
	prefix := test_utils.GenerateTempDir(t, "node_instances")
	defer os.RemoveAll(prefix)

	got := generateDirName(prefix)
	trimmed := strings.TrimPrefix(got, filepath.Join(prefix, nodeDirPrefix))

	if _, err := strconv.Atoi(trimmed); err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
	}
}
