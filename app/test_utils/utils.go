package test_utils

import (
	"bufio"
	"node_manager/app/constants"
	"os"
	"path/filepath"
	"testing"
)

const nodeCliContent = `
import time
import sys

if len(sys.argv) > 1:
    server = sys.argv[2]
    api_key = sys.argv[4]
else:
    server = input("server > ")
    api_key = input("username > ")

print("-"*10)
print("starting execution")
print("server:", server)
print("api_key:", api_key)

time.sleep(0.1)

print("stopping execution")
`

func GenerateTempDir(t testing.TB, pattern string) string {
	t.Helper()
	prefix, err := os.MkdirTemp(os.TempDir(), pattern)
	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}
	return prefix
}

func GenerateTempNode(t testing.TB, skipList []string) (nodeDir string) {
	t.Helper()
	nodeDir = GenerateTempDir(t, "zeuz_node")

	_ = os.WriteFile(filepath.Join(nodeDir, "node_cli.py"), []byte(nodeCliContent), constants.OS_USER_RW)
	tempFile, err := os.OpenFile(filepath.Join(nodeDir, "node_cli.py"), os.O_RDWR, constants.OS_USER_RW)
	//tempFile, err := ioutil.TempFile(nodeDir, "node_cli_*.py")
	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}

	writer := bufio.NewWriter(tempFile)
	_, err = writer.WriteString(nodeCliContent)
	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}

	// Write contents to disk and close the file.
	_ = writer.Flush()
	_ = tempFile.Close()

	GenerateSkipListFiles(t, nodeDir, skipList)

	return
}

func GenerateSkipListFiles(t testing.TB, nodeDir string, skipList []string) {
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
