package poll_node_state

import (
	"context"
	"fmt"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestService(t *testing.T) {
	ctx := context.Background()
	testMapFS := setupFS(t)
	srv := Service{
		fs: testMapFS,
	}

	got, _ := srv.Run(ctx, nil)
	got = got.(Result)

	want := Result{
		InProgress: 2,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func setupFS(t testing.TB) fs.FS {
	t.Helper()

	nodesWithStatus := [][]string{
		{fmt.Sprintf("node1/%s", statusFileName), "in_progress"},
		{fmt.Sprintf("node2/"), ""},
		{fmt.Sprintf("node3/"), ""},
		{fmt.Sprintf("node4/%s", statusFileName), "in_progress"},
		{fmt.Sprintf("node5/%s", statusFileName), "complete"},
		{fmt.Sprintf("node6/%s", statusFileName), "complete"},
	}

	testMapFS := fstest.MapFS{}
	for _, n := range nodesWithStatus {
		testMapFS[n[0]] = &fstest.MapFile{Data: []byte(n[1])}
	}

	noStatusNodes := []string{
		nodesWithStatus[1][0],
		nodesWithStatus[2][0],
	}

	for _, n := range noStatusNodes {
		// Set the `fs.MapFile.ModeDir` bit to mark it as a directory
		testMapFS[n].Mode = testMapFS[n].Mode | fs.ModeDir
	}

	return testMapFS
}
