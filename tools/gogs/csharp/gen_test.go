package csharp

import (
	"os"
	"testing"

	"github.com/metagogs/gogs/utils/filex"
)

func TestNewGen(t *testing.T) {
	type args struct {
		proto    string
		onlyCode bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new gen",
			args: args{
				proto:    "testdata/data.proto",
				onlyCode: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewCSharpGen(tt.args.proto, tt.args.onlyCode)
			g.Home = "test/"
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = g.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var haveErr bool
			if ok := filex.IsFileEqual("testdata/Model/Register.cs", "test/Model/Register.cs"); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", "test/Model/Register.cs", "testdata/Model/Register.cs")
				haveErr = true
			}

			if !haveErr {
				_ = os.RemoveAll("test")
			}
		})
	}
}
