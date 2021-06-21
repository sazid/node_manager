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
	fsys, validFileCount := setupFS(t)
	nodeRemover := &spyNodeRemover{}
	var srv app.Service = New(fsys, nodeRemover)

	_, err := srv.Run(context.Background(), nil)
	if err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
	}

	if nodeRemover.called != validFileCount {
		t.Errorf("got node remover called %v times, want %v", nodeRemover.called, validFileCount)
	}
}

func setupFS(t testing.TB) (fs.FS, int) {
	t.Helper()

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

	return testMapFS, len(nodePidFiles) - 2
}
