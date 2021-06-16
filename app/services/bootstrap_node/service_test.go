package bootstrap_node

import (
	"bufio"
	"context"
	"io/ioutil"
	"node_manager/app/constants"
	"node_manager/app/store"
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

	config := store.NewConfig()

	nodeDir, copyDest, err := generateTempNode(t)
	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}
	defer func() {
		_ = os.RemoveAll(nodeDir)
		_ = os.RemoveAll(copyDest)
	}()

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
	prefix := generateTempDir(t, "node_instances")
	defer os.RemoveAll(prefix)

	got := generateDirName(prefix)
	trimmed := strings.TrimPrefix(got, filepath.Join(prefix, nodeDirPrefix))

	if _, err := strconv.Atoi(trimmed); err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
	}
}

func generateTempDir(t testing.TB, pattern string) string {
	t.Helper()
	prefix, err := os.MkdirTemp(os.TempDir(), pattern)
	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}
	return prefix
}

func generateTempNode(t testing.TB) (nodeDir, copyDest string, err error) {
	t.Helper()
	nodeDir = generateTempDir(t, "zeuz_node")
	copyDest = generateTempDir(t, "node_instances")

	const fileContent = `print("Hello World")`

	tempFile, err := ioutil.TempFile(nodeDir, "node_cli_*.py")
	if err != nil {
		return
	}

	writer := bufio.NewWriter(tempFile)
	_, err = writer.WriteString(fileContent)
	if err != nil {
		return
	}

	// Write contents to disk and close the file.
	_ = writer.Flush()
	_ = tempFile.Close()

	generateSkipListFiles(t, nodeDir)

	return
}

func generateSkipListFiles(t testing.TB, nodeDir string) {
	t.Helper()
	for _, item := range skipList {
		_ = os.WriteFile(filepath.Join(nodeDir, item), []byte("testing skip files"), constants.OS_USER_RW)
		tempFile, err := os.OpenFile(filepath.Join(nodeDir, item), os.O_RDWR, constants.OS_USER_RW)
		if err != nil {
			continue
		}

		writer := bufio.NewWriter(tempFile)
		_, err = writer.WriteString("testing skip files")
		if err != nil {
			return
		}

		// Write contents to disk and close the file.
		_ = writer.Flush()
		_ = tempFile.Close()
	}
}
