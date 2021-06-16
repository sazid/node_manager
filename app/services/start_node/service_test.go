package start_node

import (
	"context"
	"testing"
)

func TestStartNode(t *testing.T) {
	ctx := context.Background()
	srv := Service{}

	_, err := srv.Run(ctx, nil)

	if err != nil {
		t.Errorf("did not expect an error, got %+v, want %+v", err, nil)
	}
}
