package bootstrap_node

import (
	"context"
	"reflect"
	"testing"
)

func TestBootstrapNode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := Service{}

	nodePath, err := srv.Run(ctx, nil)

	if err != nil {
		t.Fatalf("did not expect an error, got %+v, want %+v", err, nil)
	}

	want := Result{}
	if !reflect.DeepEqual(nodePath, want) {
		t.Errorf("got %+v, want %+v", nodePath, want)
	}
}
