package gomod

import (
	"encoding/json"

	"github.com/metagogs/gogs/utils/execx"
)

type GoModule struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

func (g *GoModule) IsInGoMod() bool {
	if len(g.Path) == 0 {
		return false
	}
	if g.Path == "command-line-arguments" {
		return false
	}
	if len(g.GoMod) == 0 {
		return false
	}

	return true
}

func GetMod() (*GoModule, error) {
	data, err := execx.Exec("go list -m -json")
	if err != nil {
		return nil, err
	}
	m := &GoModule{}
	if err := json.Unmarshal([]byte(data), m); err != nil {
		return nil, err
	}

	return m, nil
}
