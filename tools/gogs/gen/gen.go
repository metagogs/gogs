package gen

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/emicklei/proto"
	"github.com/ettle/strcase"
	"github.com/metagogs/gogs/tools/gogs/gen/gentemplate"
	"github.com/metagogs/gogs/tools/gogs/protoparse"
	"github.com/metagogs/gogs/utils/execx"
	"github.com/metagogs/gogs/utils/templatex"
	"github.com/pterm/pterm"
)

type Components struct {
	Components []*Component
}

type Component struct {
	Name        string
	Index       int
	BasePackage string
	Fields      []*Field
}

type Field struct {
	ComponentName  string
	ComponentIndex int
	BasePackage    string
	Package        string
	Name           string
	Index          int
	ServerMessage  bool
}

type Gen struct {
	proto       protoparse.Proto
	protoFile   string
	componets   *Components
	messages    map[string]protoparse.Message
	basePackage string
	Home        string
	logicPath   []string
	debugNoPb   bool
}

func NewGen(proto string, basePackage string) (*Gen, error) {
	protoPrase, err := protoparse.NewProtoParser().Parse(proto)
	if err != nil {
		return nil, err
	}
	return &Gen{
		protoFile:   proto,
		proto:       protoPrase,
		messages:    make(map[string]protoparse.Message),
		basePackage: basePackage,
	}, nil
}

func (g *Gen) Generate() error {
	g.init()

	if err := g.app(); err != nil {
		return err
	}
	if err := g.ep(); err != nil {
		return err
	}
	if err := g.svc(); err != nil {
		return err
	}
	if err := g.server(); err != nil {
		return err
	}
	if err := g.genAllLogic(); err != nil {
		return err
	}

	return nil
}

func (g *Gen) init() {
	g.componets = new(Components)
	for _, m := range g.proto.Message {
		g.messages[m.Name] = m
	}
	//find components
	for _, m := range g.proto.Message {
		if m.Comment == nil {
			continue
		}
		if protoparse.CommentsContains(m.Comment.Lines, "@gogs:Components") {
			for _, e := range m.Elements {
				if v, ok := e.(*proto.NormalField); ok {
					com := g.parseComponent(v)
					g.componets.Components = append(g.componets.Components, com)
				}
			}
		}
	}

	for _, c := range g.componets.Components {
		pterm.Success.Printf("Component: %s [%d]\n", c.Name, c.Index)
		for _, f := range c.Fields {
			pterm.Success.Printf("      Field: %s [%d][%v]\n", f.Name, f.Index, f.ServerMessage)
		}
	}
}

func (g *Gen) parseComponent(component *proto.NormalField) *Component {
	newCompoent := &Component{}
	newCompoent.Name = component.Name
	newCompoent.Index = component.Sequence
	newCompoent.BasePackage = g.basePackage

	for _, e := range g.messages[component.Name].Elements {
		if v, ok := e.(*proto.NormalField); ok {
			data := &Field{}
			data.ComponentName = component.Name
			data.ComponentIndex = component.Sequence
			data.Name = v.Name
			data.Index = v.Sequence
			data.Package = g.proto.PbPackage
			data.BasePackage = g.basePackage
			if g.messages[v.Name].Comment != nil {
				if protoparse.CommentsContains(g.messages[v.Name].Comment.Lines, "@gogs:ServerMessage") {
					data.ServerMessage = true
				}
			}

			newCompoent.Fields = append(newCompoent.Fields, data)
		}
	}

	return newCompoent

}

func (g *Gen) app() error {
	file := g.getAppFile()
	data := map[string]interface{}{}
	data["Package"] = g.proto.PbPackage
	data["BasePackage"] = g.basePackage

	if err := templatex.With("gogs").GoFmt(true).Parse(gentemplate.AppTpl).SaveTo(data, file, false); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Gen) ep() error {
	out := "."
	if len(g.Home) > 0 {
		out = g.Home
	}
	if !g.debugNoPb {
		protocCmd := fmt.Sprintf("protoc --go_out=%s %s", out, g.protoFile)
		fmt.Println(protocCmd)
		if _, err := execx.Exec(protocCmd); err != nil {
			fmt.Println(err.Error())
			pterm.Error.Println("run protoc error " + err.Error())
			return err
		}

	}

	file := g.getEPFile()
	data := map[string]interface{}{}
	data["Package"] = g.proto.PbPackage
	data["Components"] = g.componets.Components

	if err := templatex.With("gogs").GoFmt(true).Parse(gentemplate.EPTpl).SaveTo(data, file, true); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Gen) svc() error {
	file := g.getSvcFile()
	if err := templatex.With("gogs").GoFmt(true).Parse(gentemplate.SvcTpl).SaveTo(nil, file, false); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Gen) server() error {
	data := map[string]interface{}{}
	data["BasePackage"] = g.basePackage
	data["Package"] = g.proto.PbPackage
	data["Components"] = g.componets.Components

	file := g.getServerFile()
	if err := templatex.With("gogs").GoFmt(true).Parse(gentemplate.ServerTpl).SaveTo(data, file, true); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}

	return nil
}

func (g *Gen) genAllLogic() error {
	for _, c := range g.componets.Components {
		for _, f := range c.Fields {
			if !f.ServerMessage {
				if err := g.genLogic(f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (g *Gen) genLogic(com *Field) error {
	data := map[string]interface{}{}
	data["BasePackage"] = g.basePackage
	data["Package"] = g.proto.PbPackage
	data["LogicPackage"] = strings.ReplaceAll(strings.ToLower(com.ComponentName), "_", "")
	data["Name"] = com.Name
	data["SnakeName"] = strcase.ToSnake(com.Name)

	file := fmt.Sprintf("%sinternal/logic/%s/%s_logic.go", g.Home, data["LogicPackage"], strcase.ToSnake(com.Name))
	if err := templatex.With("gogs").GoFmt(true).Parse(gentemplate.LogicTpl).SaveTo(data, file, false); err != nil {
		pterm.Error.Printfln("generate file error " + file + ":" + err.Error())
		return err
	}
	g.logicPath = append(g.logicPath, file)

	return nil
}

func (g *Gen) getAppFile() string {
	return g.Home + "main.go"
}

func (g *Gen) getEPFile() string {
	return g.Home + fmt.Sprintf("%s/%s.ep.go", g.proto.PbPackage, path.Base(strings.ReplaceAll(g.protoFile, ".proto", "")))
}

func (g *Gen) getPBFile() string {
	return g.Home + fmt.Sprintf("%s/%s.pb.go", g.proto.PbPackage, path.Base(g.protoFile))
}

func (g *Gen) getSvcFile() string {
	return g.Home + "internal/svc/service_context.go"
}

func (g *Gen) getServerFile() string {
	return g.Home + "internal/server/server.go"
}

func (g *Gen) getLogicFile() []string {
	return g.logicPath
}

func (g *Gen) clean() {
	_ = os.Remove(g.getAppFile())
	_ = os.Remove(g.getEPFile())
	_ = os.Remove(g.getPBFile())
	_ = os.Remove(g.getSvcFile())
	_ = os.Remove(g.getServerFile())
	_ = os.Remove("internal/logic")
	_ = os.RemoveAll(g.proto.PbPackage)
}
