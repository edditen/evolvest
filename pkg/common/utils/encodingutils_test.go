package utils

import (
	"bytes"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not nil",
			args: args{
				data: bytes.Repeat([]byte{109, 45, 10, 211, 8}, 3),
			},
			want: "bS0K0whtLQrTCG0tCtMI",
		},
		{
			name: "nil",
			args: args{
				data: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64Encode(tt.args.data); got != tt.want {
				t.Errorf("Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
