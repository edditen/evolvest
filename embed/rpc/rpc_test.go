package rpc

import (
	"context"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestGetEvolvestServer(t *testing.T) {
	tests := []struct {
		name string
		want *EvolvestServer
	}{
		{
			name: "normal",
			want: NewEvolvestServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEvolvestServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvolvestServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEvolvestServer(t *testing.T) {
	tests := []struct {
		name string
		want *EvolvestServer
	}{
		{
			name: "normal",
			want: &EvolvestServer{
				store: store.GetStore(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvolvestServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvolvestServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvolvestServer_Keys(t *testing.T) {
	t.Run("pattern all", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.KeysRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.KeysResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.KeysRequest{
					Pattern: ".*",
				},
			},
			want: &evolvest.KeysResponse{
				Keys: []string{"hello", "abc"},
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Keys().Return([]string{"hello", "abc"}, nil).Times(1)

		got, err := e.Keys(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Keys() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Keys() got = %v, want %v", got, tt.want)
		}
	})

	t.Run("pattern part", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.KeysRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.KeysResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.KeysRequest{
					Pattern: "he.*",
				},
			},
			want: &evolvest.KeysResponse{
				Keys: []string{"hello", "hero"},
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Keys().Return([]string{"hello", "abc", "", "hero"}, nil).Times(1)

		got, err := e.Keys(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Keys() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Keys() got = %v, want %v", got, tt.want)
		}
	})

	t.Run("empty", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.KeysRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.KeysResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.KeysRequest{
					Pattern: "he.*",
				},
			},
			want: &evolvest.KeysResponse{
				Keys: []string{},
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Keys().Return([]string{"abc", "", "xyz"}, nil).Times(1)

		got, err := e.Keys(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Keys() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Keys() got = %v, want %v", got, tt.want)
		}
	})

}

func Test_parseCmd(t *testing.T) {
	type args struct {
		cmdText string
	}
	tests := []struct {
		name string
		args args
		want *common.TxRequest
	}{
		{
			name: "sync",
			args: args{
				cmdText: "123  sync set hello d29ybGQ=",
			},
			want: &common.TxRequest{
				TxId:   int64(123),
				Flag:   common.FlagSync,
				Action: common.SET,
				Key:    "hello",
				Val:    []byte("world"),
			},
		},
		{
			name: "req",
			args: args{
				cmdText: "123 req set hello d29ybGQ=",
			},
			want: &common.TxRequest{
				TxId:   int64(123),
				Flag:   common.FlagSync,
				Action: common.SET,
				Key:    "hello",
				Val:    []byte("world"),
			},
		},
		{
			name: "non val",
			args: args{
				cmdText: "123 sync set hello",
			},
			want: &common.TxRequest{
				TxId:   int64(123),
				Flag:   common.FlagSync,
				Action: common.SET,
				Key:    "hello",
				Val:    []byte{},
			},
		},
		{
			name: "multiple values",
			args: args{
				cmdText: "123 sync set hello world world2",
			},
			want: nil,
		},
		{
			name: "wrong format",
			args: args{
				cmdText: "123 sync set hello world",
			},
			want: nil,
		},
		{
			name: "missing required",
			args: args{
				cmdText: "123 sync set",
			},
			want: nil,
		},
		{
			name: "del",
			args: args{
				cmdText: "123 sync del hello ",
			},
			want: &common.TxRequest{
				TxId:   int64(123),
				Flag:   common.FlagSync,
				Action: common.DEL,
				Key:    "hello",
				Val:    nil,
			},
		},
		{
			name: "command not support",
			args: args{
				cmdText: "123 sync get hello ",
			},
			want: nil,
		},
		{
			name: "can not convert",
			args: args{
				cmdText: "abc sync del hello ",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCmd(tt.args.cmdText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
