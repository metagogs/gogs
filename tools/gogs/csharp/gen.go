package csharp

import (
	"fmt"
	"os"

	"github.com/emicklei/proto"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/tools/gogs/csharp/gentemplate"
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
	Action10       string //action 10进制
	Action16       string //action 16进制
}

type CSharpGen struct {
	proto       protoparse.Proto
	protoFile   string
	componets   *Components
	messages    map[string]protoparse.Message
	basePackage string
	Home        string
	logicPath   []string
	debugNoPb   bool
	onlyCode    bool //just generate code by proto exinclude gogs
}

func NewCSharpGen(proto string, onlyCode bool) (*CSharpGen, error) {
	protoPrase, err := protoparse.NewProtoParser().Parse(proto)
	if err != nil {
		return nil, err
	}
	return &CSharpGen{
		protoFile: proto,
		proto:     protoPrase,
		messages:  make(map[string]protoparse.Message),
		onlyCode:  onlyCode,
	}, nil
}

func (g *CSharpGen) Generate() error {
	g.init()

	if g.onlyCode {
		if err := g.gogs(); err != nil {
			return err
		}
	}

	if err := g.register(); err != nil {
		return err
	}

	return nil
}

func (g *CSharpGen) init() {
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

func (g *CSharpGen) parseComponent(component *proto.NormalField) *Component {
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
			//create action
			actionValue := packet.CreateAction(packet.ServicePacket, uint8(data.ComponentIndex), uint16(data.Index))
			data.Action10 = fmt.Sprint(actionValue)
			data.Action16 = fmt.Sprintf("0x%x", actionValue)
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

func (g *CSharpGen) gogs() error {
	if err := templatex.With("gogs").Parse(gentemplate.CodecTpl).SaveTo(nil, g.Home+"Gogs/Codec.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Gogs/Codec.cs:" + err.Error())
		return err
	}
	if err := templatex.With("gogs").Parse(gentemplate.CommonTpl).SaveTo(nil, g.Home+"Gogs/Common.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Gogs/Common.cs:" + err.Error())
		return err
	}
	if err := templatex.With("gogs").Parse(gentemplate.EventsManagerTpl).SaveTo(nil, g.Home+"Gogs/EventsManager.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Gogs/EventsManager.cs:" + err.Error())
		return err
	}
	if err := templatex.With("gogs").Parse(gentemplate.ICodecTpl).SaveTo(nil, g.Home+"Gogs/ICodec.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Gogs/ICodec.cs:" + err.Error())
		return err
	}
	if err := templatex.With("gogs").Parse(gentemplate.MessagesTpl).SaveTo(nil, g.Home+"Gogs/Messages.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Gogs/Messages.cs:" + err.Error())
		return err
	}
	if err := templatex.With("gogs").Parse(gentemplate.PacketTpl).SaveTo(nil, g.Home+"Gogs/Packet.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Gogs/Packet.cs:" + err.Error())
		return err
	}

	return nil
}

func (g *CSharpGen) register() error {
	out := "Model"
	if len(g.Home) > 0 {
		out = g.Home + "Model"
	}
	if !g.debugNoPb {
		os.MkdirAll(out, os.ModePerm)
		protocCmd := fmt.Sprintf("protoc --csharp_out=%s %s", out, g.protoFile)
		if _, err := execx.Exec(protocCmd); err != nil {
			fmt.Println(err.Error())
			pterm.Error.Println("run protoc error " + err.Error())
			return err
		}

	}

	data := map[string]interface{}{}
	data["Package"] = g.proto.PbPackage
	data["Components"] = g.componets.Components

	if err := templatex.With("gogs").Parse(gentemplate.RegisterTpl).SaveTo(data, g.Home+"Model/Register.cs", true); err != nil {
		pterm.Error.Printfln("generate file error Model/Register.cs:" + err.Error())
		return err
	}

	return nil
}
