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

func TestServiceCancelled(t *testing.T) {
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
		{fmt.Sprintf("node1/%s", statusFileName), statusInProgress},
		{fmt.Sprintf("node2/%s", statusFileName), statusIdle},
		{fmt.Sprintf("node3/%s", statusFileName), statusIdle},
		{fmt.Sprintf("node4/%s", statusFileName), statusInProgress},
		{fmt.Sprintf("node5/%s", statusFileName), statusComplete},
		{fmt.Sprintf("node6/%s", statusFileName), statusComplete},
		{fmt.Sprintf("node7/"), statusComplete},
		{fmt.Sprintf("/"), ""},
	}

	testMapFS := fstest.MapFS{}
	for _, n := range nodesWithStatus {
		testMapFS[n[0]] = &fstest.MapFile{Data: []byte(n[1])}
	}

	return testMapFS
}
