package gen

import (
	"os"
	"testing"

	"github.com/metagogs/gogs/utils/filex"
)

func TestNewGen(t *testing.T) {
	type args struct {
		proto       string
		basePackage string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new gen",
			args: args{
				proto:       "testdata/data.proto",
				basePackage: "github.com/metagogs/test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGen(tt.args.proto, tt.args.basePackage)
			g.Home = "test/"
			g.debugNoPb = true
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = g.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var haveErr bool
			if ok := filex.IsFileEqual("testdata/main.go", g.getAppFile()); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getAppFile(),
					"testdata/main.go")
				haveErr = true
			}
			if ok := filex.IsFileEqual("testdata/model/data.ep.go", g.getEPFile()); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getEPFile(),
					"testdata/data.ep.go")
				haveErr = true
			}
			if ok := filex.IsFileEqual("testdata/internal/svc/service_context.go", g.getSvcFile()); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getSvcFile(),
					"testdata/internal/svc/service_context.go")
				haveErr = true
			}
			if ok := filex.IsFileEqual("testdata/internal/server/server.go", g.getServerFile()); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getServerFile(),
					"testdata/internal/server/server.go")
				haveErr = true
			}
			if ok := filex.IsFileEqual("testdata/internal/logic/baseworld/bind_user_logic.go", g.getLogicFile()[0]); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", g.getLogicFile()[0],
					"testdata/internal/logic/baseworld/bind_user_logic.go")
				haveErr = true
			}

			if !haveErr {
				g.clean()
				_ = os.RemoveAll("test")
			}
		})
	}
}
