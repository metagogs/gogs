package execx

import (
	"fmt"
	"runtime"
	"testing"
)

func TestExec(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "go version",
			args: args{
				arg: "go version",
			},
			want: fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		},
		{
			name: "error command",
			args: args{
				arg: "go ver",
			},
			wantErr: true,
		},
		{
			name: "error command",
			args: args{
				arg: "ggggg",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exec(tt.args.arg)
			fmt.Print("get", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
