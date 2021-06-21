package poll_node_state

import (
	"context"
	"fmt"
	"io/fs"
	"node_manager/app"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestPollService(t *testing.T) {
	ctx := context.Background()
	fsys := setupFS(t)
	srv := Service{
		fsys: fsys,
	}

	got, _ := srv.Run(ctx, nil)
	got = got.(Result)

	want := Result{
		Complete:   2,
		Idle:       2,
		InProgress: 2,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestPollServiceCanBeCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fsys := setupFS(t)
	srv := Service{
		fsys: fsys,
	}

	cancel()
	_, err := srv.Run(ctx, nil)

	if err != ErrServiceCancelled {
		t.Errorf("expected the service to be cancelled, got %+v, want %+v", err, ErrServiceCancelled)
	}
}

func setupFS(t testing.TB) fs.FS {
	t.Helper()

	nodesWithStatus := [][]string{
		{fmt.Sprintf("node1/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateInProgress)},
		{fmt.Sprintf("node2/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateIdle)},
		{fmt.Sprintf("node3/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateIdle)},
		{fmt.Sprintf("node4/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateInProgress)},
		{fmt.Sprintf("node5/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateComplete)},
		{fmt.Sprintf("node6/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateComplete)},
		{fmt.Sprintf("node7/"), ""}, // no `node_state.json` file
		{fmt.Sprintf("/"), ""},      // invalid path
	}

	testMapFS := fstest.MapFS{}
	for _, n := range nodesWithStatus {
		testMapFS[n[0]] = &fstest.MapFile{Data: []byte(n[1])}
	}

	return testMapFS
}
