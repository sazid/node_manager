package node_remover

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestRemoveNode(t *testing.T) {
	srv := &Service{}

	tempDir, err := ioutil.TempDir(os.TempDir(), "node_*")
	if err != nil {
		t.Fatalf("should not have failed to create temp dir, got %+v", err)
	}

	fmt.Println(tempDir)

	msg := Message{
		Dir: tempDir,
	}
	_, err = srv.Run(context.Background(), msg)

	if err != nil {
		t.Fatalf("did not expect an error, got %+v", err)
	}

	_, err = os.Stat(tempDir)
	if !os.IsNotExist(err) {
		t.Fatal("service did not remove the specified node directory.", err)
	}
}
