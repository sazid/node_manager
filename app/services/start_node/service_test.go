package start_node

import (
	"bytes"
	"context"
	"io"
	"log"
	"node_manager/app"
	"node_manager/app/services/bootstrap_node"
	"node_manager/app/store"
	"node_manager/app/test_utils"
	"os"
	"testing"
)

func TestStartNode(t *testing.T) {
	config := store.DummyConfig(t, 1, 1,
		"https://github.com",
		"123ABC-456DEF-789GHI-101JKL")
	ctx := context.Background()

	buf := &bytes.Buffer{}

	srv := Service{
		Config:           config,
		BootstrapNodeSrv: stubBootstrapNodeSrv(t),
		OutputWriter:     buf,
	}

	type tempResult struct {
		res Result
		err error
	}
	ch := make(chan tempResult)
	go func() {
		res, err := srv.Run(ctx, nil)
		ch <- tempResult{res.(Result), err}
	}()

	result := <-ch
	_ = os.RemoveAll(result.res.NodePath)

	log.Println("Output:")
	_, _ = io.Copy(os.Stdout, buf)

	if result.err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", result.err, nil)
	}
}

func stubBootstrapNodeSrv(t testing.TB) app.ServiceFunc {
	return func(context.Context, interface{}) (result interface{}, err error) {
		nodeDir := test_utils.GenerateTempNode(t, nil)

		return bootstrap_node.Result{
			Path: nodeDir,
		}, nil
	}
}
