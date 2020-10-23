package store

import (
	"reflect"
	"testing"
)

func TestEvolvest_Del(t *testing.T) {
	type fields struct {
		storage map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal interface{}
		wantErr bool
	}{
		{
			name: "key exits",
			fields: fields{
				storage: map[string]interface{}{
					"hello": "world",
				},
			},
			args: args{
				key: "hello",
			},
			wantVal: "world",
			wantErr: false,
		},
		{
			name: "key not exits",
			fields: fields{
				storage: map[string]interface{}{
					"hello": "world",
				},
			},
			args: args{
				key: "hello123",
			},
			wantVal: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Evolvest{
				storage: tt.fields.storage,
			}
			gotVal, err := e.Del(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Del() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestEvolvest_Get(t *testing.T) {
	type fields struct {
		storage map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal interface{}
		wantErr bool
	}{
		{
			name: "key exits",
			fields: fields{
				storage: map[string]interface{}{
					"hello": "world",
				},
			},
			args: args{
				key: "hello",
			},
			wantVal: "world",
			wantErr: false,
		},
		{
			name: "key not exits",
			fields: fields{
				storage: map[string]interface{}{
					"hello": "world",
				},
			},
			args: args{
				key: "hello123",
			},
			wantVal: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Evolvest{
				storage: tt.fields.storage,
			}
			gotVal, err := e.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Get() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestEvolvest_Set(t *testing.T) {
	type fields struct {
		storage map[string]interface{}
	}
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOldVal interface{}
		wantErr    bool
	}{
		{
			name: "key exit",
			fields: fields{
				storage: map[string]interface{}{
					"hello": "world",
				},
			},
			args: args{
				key: "hello",
				val: "123",
			},
			wantOldVal: "world",
			wantErr:    false,
		},
		{
			name: "key not exit",
			fields: fields{
				storage: map[string]interface{}{
					"hello": "world",
				},
			},
			args: args{
				key: "hello123",
				val: "123",
			},
			wantOldVal: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Evolvest{
				storage: tt.fields.storage,
			}
			gotOldVal, err := e.Set(tt.args.key, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOldVal, tt.wantOldVal) {
				t.Errorf("Set() gotOldVal = %v, want %v", gotOldVal, tt.wantOldVal)
			}
		})
	}
}

func TestNewEvolvest(t *testing.T) {
	tests := []struct {
		name string
		want *Evolvest
	}{
		{
			name: "normal",
			want: &Evolvest{storage: make(map[string]interface{}, 17)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvolvest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvolvest() = %v, want %v", got, tt.want)
			}
		})
	}
}
