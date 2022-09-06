package gen

import (
	"os"
	"runtime"
	"strings"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/tools/gogs/gen/gentemplate"
	"github.com/metagogs/gogs/utils/templatex"
	"github.com/pterm/pterm"
)

type Init struct {
	PackageName string
	Home        string
}

func NewInit(name string) *Init {
	return &Init{
		PackageName: name,
	}
}

func (g *Init) Generate() error {
	if err := g.goMod(); err != nil {
		return err
	}
	if err := g.config(); err != nil {
		return err
	}
	if err := g.proto(); err != nil {
		return err
	}

	return nil
}

func (g *Init) goMod() error {
	file := g.getGoModFile()
	data := map[string]string{
		"ProjectPackage": g.PackageName,
		"GoVersion":      strings.ReplaceAll(runtime.Version(), "go", "go "),
		"GoGSVersion":    gogs.Version,
	}
	if err := templatex.With("gogs").GoFmt(false).Parse(gentemplate.GoModTpl).SaveTo(data, file, false); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Init) config() error {
	file := g.getConfigFile()
	if err := templatex.With("gogs").GoFmt(false).Parse(gentemplate.ConfigTpl).SaveTo(nil, file, false); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Init) proto() error {
	file := g.getProtoFile()
	if err := templatex.With("gogs").GoFmt(false).Parse(gentemplate.ProtoTpl).SaveTo(nil, file, false); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Init) getGoModFile() string {
	return g.Home + "go.mod"
}

func (g *Init) getConfigFile() string {
	return g.Home + "config.yaml"
}

func (g *Init) getProtoFile() string {
	return g.Home + "data.proto"
}

func (g *Init) clean() {
	_ = os.Remove(g.getGoModFile())
	_ = os.Remove(g.getConfigFile())
	_ = os.Remove(g.getProtoFile())
}
