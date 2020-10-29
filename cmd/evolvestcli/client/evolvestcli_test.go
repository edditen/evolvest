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

		err := GetEvolvestClient().Set(context.Background(), "hello", "123")
		if err != nil {
			t.Errorf("set val error, %v\n", err)
		}

		err = GetEvolvestClient().Set(context.Background(), "hello", "world")
		if err != nil {
			t.Errorf("set val error, %v\n", err)
		}
		val, err := GetEvolvestClient().Get(context.Background(), "hello")
		if err != nil {
			t.Errorf("get val error, %v\n", err)
		}

		if val != "world" {
			t.Errorf("got val: %s, want val: %s", val, "world")
		}
		err = GetEvolvestClient().Del(context.Background(), "hello")
		if err != nil {
			t.Errorf("del val error, %v\n", err)
		}

		err = GetEvolvestClient().Del(context.Background(), "hello")
		t.Log("duplicated del, ", err)
		if err == nil {
			t.Errorf("twice del val error, %v\n", err)
		}
	})
}
