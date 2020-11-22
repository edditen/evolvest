package store

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEvolvest_Del(t *testing.T) {
	type fields struct {
		Nodes map[string]ValItem
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal ValItem
		wantErr bool
	}{
		{
			name: "key exits",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			args: args{
				key: "hello",
			},
			wantVal: ValItem{
				Val:     []byte("world"),
				Version: 123,
			},
			wantErr: false,
		},
		{
			name: "key not exits",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			args: args{
				key: "hello123",
			},
			wantVal: ValItem{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			e := &Evolvest{
				Nodes: tt.fields.Nodes,
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
		Nodes map[string]ValItem
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal ValItem
		wantErr bool
	}{
		{
			name: "key exits",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			args: args{
				key: "hello",
			},
			wantVal: ValItem{
				Val:     []byte("world"),
				Version: 123,
			},
			wantErr: false,
		},
		{
			name: "key not exits",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			args: args{
				key: "hello123",
			},
			wantVal: ValItem{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			e := &Evolvest{
				Nodes: tt.fields.Nodes,
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
		Nodes map[string]ValItem
	}
	type args struct {
		key string
		val ValItem
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOldVal ValItem
		wantExist  bool
	}{
		{
			name: "key exit",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			args: args{
				key: "hello",
				val: ValItem{
					Val:     []byte("123"),
					Version: 123,
				},
			},
			wantOldVal: ValItem{
				Val:     []byte("world"),
				Version: 123,
			},
			wantExist: true,
		},
		{
			name: "key not exit",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			args: args{
				key: "hello123",
				val: ValItem{
					Val:     []byte("123"),
					Version: 123,
				},
			},
			wantOldVal: ValItem{},
			wantExist:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			e := &Evolvest{
				Nodes: tt.fields.Nodes,
			}
			gotOldVal, goExist := e.Set(tt.args.key, tt.args.val)
			if goExist != tt.wantExist {
				t.Errorf("Set() exit = %v, wantExist %v", goExist, tt.wantExist)
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
			want: &Evolvest{Nodes: make(map[string]ValItem, 17)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			if got := NewEvolvest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvolvest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvolvest_Save(t *testing.T) {
	type fields struct {
		Nodes map[string]ValItem
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		{
			name: "nil",
			fields: fields{
				Nodes: nil,
			},
			wantData: []byte(`{"nodes":null}`),
			wantErr:  false,
		},
		{
			name: "empty",
			fields: fields{
				Nodes: map[string]ValItem{},
			},
			wantData: []byte(`{"nodes":{}}`),
			wantErr:  false,
		},
		{
			name: "have values",
			fields: fields{
				Nodes: map[string]ValItem{
					"hello": {
						Val:     []byte("world"),
						Version: 123,
					},
				},
			},
			wantData: []byte(`{"nodes":{"hello":{"Val":"d29ybGQ=","Version":123}}}`),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			e := &Evolvest{
				Nodes: tt.fields.Nodes,
			}
			gotData, err := e.Serialize()
			fmt.Printf("gotData: %s\n", string(gotData))
			if (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Serialize() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestEvolvest_Load(t *testing.T) {
	type fields struct {
		Nodes map[string]ValItem
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Nodes: map[string]ValItem{
					"abc": {
						Val:     []byte("123"),
						Version: 123,
					},
				},
			},
			args: args{
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "empty",
			fields: fields{
				Nodes: map[string]ValItem{
					"abc": {
						Val:     []byte("123"),
						Version: 123,
					},
				},
			},
			args: args{
				data: []byte("{}"),
			},
			wantErr: false,
		},
		{
			name: "not empty",
			fields: fields{
				Nodes: map[string]ValItem{
					"abc": {
						Val:     []byte("123"),
						Version: 123,
					},
					"hello": {
						Val:     []byte("456"),
						Version: 123,
					},
				},
			},
			args: args{
				data: []byte(`{"nodes":{"abc":{"Val":"MTIz","Version":123},"hello":{"Val":null,"Version":0}}}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			e := &Evolvest{
				Nodes: tt.fields.Nodes,
			}
			if err := e.Load(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			d, err := e.Serialize()
			fmt.Printf("data:%s, err:%v\n", string(d), err)
		})
	}
}
