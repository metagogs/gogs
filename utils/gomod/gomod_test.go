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
