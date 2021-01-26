package dns

import (
	"reflect"
	"testing"
)

func TestIP2A(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    [4]byte
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				ip: "1.1.1.1",
			},
			want:    [4]byte{1, 1, 1, 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IP2A(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IP2A() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IP2A() got = %v, want %v", got, tt.want)
			}
		})
	}
}
