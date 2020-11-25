package client

import (
	"context"
	"testing"
)

func TestEvolvestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	t.Run("client", func(t *testing.T) {
		port := ":8762"
		StartClient(port)

		_, err := GetEvolvestClient().Keys(context.Background(), ".*")
		if err != nil {
			t.Errorf("keys error, %v\n", err)
		}

	})
}
