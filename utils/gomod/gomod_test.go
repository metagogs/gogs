package gomod

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestGetMod(t *testing.T) {
	tests := []struct {
		name    string
		want    *GoModule
		wantErr bool
	}{
		{
			name: "get go mod",
			want: &GoModule{
				Path:      "github.com/metagogs/gogs",
				Main:      true,
				GoVersion: strings.Replace(runtime.Version(), "go", "", 1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMod()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Path, tt.want.Path) {
				t.Errorf("GetMod() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Main, tt.want.Main) {
				t.Errorf("GetMod() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.GoVersion, tt.want.GoVersion) {
				t.Errorf("GetMod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoModule_IsInGoMod(t *testing.T) {
	type fields struct {
		Path      string
		Main      bool
		Dir       string
		GoMod     string
		GoVersion string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "check is in go mod",
			fields: fields{
				Path:      "github.com/metagogs/gogs",
				Main:      true,
				Dir:       "/Users/neo/Documents/gogs",
				GoMod:     "/Users/neo/Documents/gogs/go.mod",
				GoVersion: "1.19",
			},
			want: true,
		},
		{
			name: "check is not in go mod",
			fields: fields{
				Path:      "command-line-arguments",
				Main:      true,
				Dir:       "/Users/neo/Documents/gogs",
				GoMod:     "/Users/neo/Documents/gogs/go.mod",
				GoVersion: "1.19",
			},
			want: false,
		},
		{
			name: "check is not in go mod",
			fields: fields{
				Path:      "",
				Main:      true,
				Dir:       "/Users/neo/Documents/gogs",
				GoMod:     "/Users/neo/Documents/gogs/go.mod",
				GoVersion: "1.19",
			},
			want: false,
		},
		{
			name: "check is not in go mod",
			fields: fields{
				Path:      "github.com/metagogs/gogs",
				Main:      true,
				Dir:       "/Users/neo/Documents/gogs",
				GoMod:     "",
				GoVersion: "1.19",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoModule{
				Path:      tt.fields.Path,
				Main:      tt.fields.Main,
				Dir:       tt.fields.Dir,
				GoMod:     tt.fields.GoMod,
				GoVersion: tt.fields.GoVersion,
			}
			if got := g.IsInGoMod(); got != tt.want {
				t.Errorf("GoModule.IsInGoMod() = %v, want %v", got, tt.want)
			}
		})
	}
}
