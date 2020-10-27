package rpc

import (
	"context"
	"fmt"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestEvolvestServer_Del(t *testing.T) {

	t.Run("key exist", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.DelRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.DelResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.DelRequest{
					Key: "hello",
				},
			},
			want: &evolvest.DelResponse{
				Key: "hello",
				Val: "world",
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Del(gomock.Any()).Return("world", nil).Times(1)

		got, err := e.Del(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Del() got = %v, want %v", got, tt.want)
		}
	})

	t.Run("key not exist", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.DelRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.DelResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.DelRequest{
					Key: "hello",
				},
			},
			want:    nil,
			wantErr: true,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Del(gomock.Any()).Return("", fmt.Errorf("key not exist")).Times(1)

		got, err := e.Del(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Del() got = %v, want %v", got, tt.want)
		}
	})
}

func TestEvolvestServer_Get(t *testing.T) {
	t.Run("key exist", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.GetRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.GetResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.GetRequest{
					Key: "hello",
				},
			},
			want: &evolvest.GetResponse{
				Key: "hello",
				Val: "world",
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Get(gomock.Any()).Return("world", nil).Times(1)

		got, err := e.Get(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Get() got = %v, want %v", got, tt.want)
		}
	})

	t.Run("key not exist", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.GetRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.GetResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.GetRequest{
					Key: "hello",
				},
			},
			want:    nil,
			wantErr: true,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Get(gomock.Any()).Return("", fmt.Errorf("key not exist")).Times(1)

		got, err := e.Get(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Del() got = %v, want %v", got, tt.want)
		}
	})
}

func TestEvolvestServer_Set(t *testing.T) {
	t.Run("key not exist", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.SetRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.SetResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.SetRequest{
					Key: "hello",
					Val: "world",
				},
			},
			want: &evolvest.SetResponse{
				Key:      "hello",
				ExistVal: false,
				OldVal:   "",
				NewVal:   "world",
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return("", false).Times(1)

		got, err := e.Set(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Set() got = %v, want %v", got, tt.want)
		}
	})

	t.Run("key exist", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()
		mockStore := store.NewMockStore(controller)

		type fields struct {
			store store.Store
		}
		type args struct {
			ctx     context.Context
			request *evolvest.SetRequest
		}
		tt := struct {
			fields  fields
			args    args
			want    *evolvest.SetResponse
			wantErr bool
		}{
			fields: fields{
				store: mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &evolvest.SetRequest{
					Key: "hello",
					Val: "world",
				},
			},
			want: &evolvest.SetResponse{
				Key:      "hello",
				ExistVal: true,
				OldVal:   "123",
				NewVal:   "world",
			},
			wantErr: false,
		}

		e := &EvolvestServer{
			store: tt.fields.store,
		}

		mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return("123", true).Times(1)

		got, err := e.Set(tt.args.ctx, tt.args.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Set() got = %v, want %v", got, tt.want)
		}
	})

}

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
