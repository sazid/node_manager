package kill_node

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"node_manager/app"
	"node_manager/app/services/node_remover"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
)

type spyNodeRemover struct {
	called int
}

func (s *spyNodeRemover) Run(_ context.Context, message interface{}) (result interface{}, err error) {
	m := message.(node_remover.Message)
	if !strings.Contains(m.NodeAbsolutePath, app.PidFilename) {
		return nil, errors.New("does not contain the PID file")
	}
	s.called++
	return
}

func TestKillNode(t *testing.T) {
	fsys, idleNodeCount := setupFS(t)
	nodeRemover := &spyNodeRemover{}
	var srv app.Service = New(fsys, nodeRemover)

	_, err := srv.Run(context.Background(), nil)
	if err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
	}

	if nodeRemover.called != idleNodeCount {
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
		{fmt.Sprintf("node7/"), ""}, // no `node_state.json` file
		{fmt.Sprintf("/"), ""},      // invalid path
	}

	nodePidFiles := [][]string{
		{fmt.Sprintf("node1/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node2/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node3/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node4/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node5/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node6/%s", app.PidFilename), strconv.Itoa(app.PidSentinelValue)},
		{fmt.Sprintf("node7/"), ""}, // no `pid.txt` file
		{fmt.Sprintf("/"), ""},      // invalid path
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
