package kill_node

import (
	"context"
	"fmt"
	"io/fs"
	"node_manager/app"
	"node_manager/app/services/node_remover"
	"strconv"
	"testing"
	"testing/fstest"
)

type spyNodeRemover struct {
	called int
}

func (s *spyNodeRemover) Run(_ context.Context, message interface{}) (result interface{}, err error) {
	if _, ok := message.(node_remover.Message); !ok {
		return nil, fmt.Errorf("`message` type should be `node_remover.Message`, got %T", message)
	}
	s.called++
	return
}

func TestKillNode(t *testing.T) {
	fsys, idleNodeCount := setupFS(t)
	nodeRemover := &spyNodeRemover{}
	var srv app.Service = New(fsys, ".", nodeRemover)

	_, err := srv.Run(context.Background(), nil)
	if err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
	}

	// TODO (report): remove the +1 from idleNodeCount once the report service
	// added.
	if nodeRemover.called+1 != idleNodeCount {
		t.Errorf("got node remover called %v times, want %v", nodeRemover.called, idleNodeCount)
	}
}

func setupFS(t testing.TB) (fsys fs.FS, idleNodeCount int) {
	t.Helper()

	nodesWithStatus := [][]string{
		{fmt.Sprintf("node1/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateInProgress)},
		{fmt.Sprintf("node2/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateIdle)},
		{fmt.Sprintf("node3/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateIdle)},
		{fmt.Sprintf("node4/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateInProgress)},
		{fmt.Sprintf("node5/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateComplete)},
		{fmt.Sprintf("node6/%s", app.NodeStateFilename), fmt.Sprintf(app.StatusTemplate, app.StateComplete)},
		{"node7/", ""}, // no `node_state.json` file
		{"/", ""},      // invalid path
	}

	nodePidFiles := [][]string{
		{fmt.Sprintf("node1/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node2/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node3/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node4/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node5/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node6/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{"node7/", ""}, // no `pid.txt` file
		{"/", ""},      // invalid path
	}

	testMapFS := fstest.MapFS{}
	for _, n := range nodePidFiles {
		testMapFS[n[0]] = &fstest.MapFile{Data: []byte(n[1])}
	}
	for _, n := range nodesWithStatus {
		testMapFS[n[0]] = &fstest.MapFile{Data: []byte(n[1])}
	}

	return testMapFS, 2
}
