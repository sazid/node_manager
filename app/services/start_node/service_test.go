package start_node

import (
	"context"
	"node_manager/app"
	"node_manager/app/services/bootstrap_node"
	"node_manager/app/test_utils"
	"os"
	"testing"
)

func TestStartNode(t *testing.T) {
	ctx := context.Background()

	srv := Service{
		bootstrapNodeSrv: stubBootstrapNodeSrv(t),
	}

	res, err := srv.Run(ctx, nil)
	defer os.RemoveAll(res.(Result).Path)

	if err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
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
