package gen

import (
	"os"
	"testing"

	"github.com/metagogs/gogs/utils/filex"
)

func TestInit_Generate(t *testing.T) {
	type fields struct {
		PackageName string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "generate",
			fields: fields{
				PackageName: "github.com/metagogs/test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Init{
				PackageName: tt.fields.PackageName,
			}
			g.Home = "test/"
			if err := g.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("Init.Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			var haveErr bool

			if ok := filex.IsFileEqual("testdata/config.yaml", g.getConfigFile()); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getGoModFile(), "testdata/config.yaml")
				haveErr = true
			}
			if ok := filex.IsFileEqual("testdata/data.proto", g.getProtoFile()); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getProtoFile(), "testdata/data.proto")
				haveErr = true
			}

			if !haveErr {
				g.clean()
				_ = os.Remove("test")
			}

		})
	}
}
