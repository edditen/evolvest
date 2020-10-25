package client

import (
	"context"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"reflect"
	"testing"
)

func TestEvolvestClient_Del(t *testing.T) {
	type fields struct {
		client evolvest.EvolvestServiceClient
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EvolvestClient{
				client: tt.fields.client,
			}
			if err := e.Del(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEvolvestClient_Get(t *testing.T) {
	type fields struct {
		client evolvest.EvolvestServiceClient
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EvolvestClient{
				client: tt.fields.client,
			}
			gotVal, err := e.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Get() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestEvolvestClient_Set(t *testing.T) {
	type fields struct {
		client evolvest.EvolvestServiceClient
	}
	type args struct {
		ctx context.Context
		key string
		val string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EvolvestClient{
				client: tt.fields.client,
			}
			if err := e.Set(tt.args.ctx, tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEvolvestClient(t *testing.T) {
	tests := []struct {
		name string
		want *EvolvestClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEvolvestClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvolvestClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEvolvestClient(t *testing.T) {
	tests := []struct {
		name string
		want *EvolvestClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvolvestClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvolvestClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartClient(t *testing.T) {
	type args struct {
		port string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestEvolvestClient(t *testing.T) {
	t.Run("client", func(t *testing.T) {
		port := ":8762"
		StartClient(port)
		err := GetEvolvestClient().Set(context.Background(), "hello", "world")
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
	})
}
