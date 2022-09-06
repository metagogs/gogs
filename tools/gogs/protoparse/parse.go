package protoparse

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/emicklei/proto"
)

type Package struct {
	*proto.Package
}

type Import struct {
	*proto.Import
}

type Message struct {
	*proto.Message
}

type Service struct {
	*proto.Service
	RPC []*RPC
}

type RPC struct {
	*proto.RPC
}

type Proto struct {
	Src       string
	Name      string
	Package   Package
	PbPackage string
	GoPackage string
	Import    []Import
	Message   []Message
	Service   Service
}
type ProtoParser struct{}

func NewProtoParser() *ProtoParser {
	return &ProtoParser{}
}

func (p *ProtoParser) Parse(src string) (Proto, error) {
	var ret Proto

	abs, err := filepath.Abs(src)
	if err != nil {
		return Proto{}, err
	}

	r, err := os.Open(abs)
	if err != nil {
		return ret, err
	}
	defer r.Close()

	parser := proto.NewParser(r)
	set, err := parser.Parse()
	if err != nil {
		return ret, err
	}

	var serviceList []Service
	proto.Walk(
		set,
		proto.WithImport(func(i *proto.Import) {
			ret.Import = append(ret.Import, Import{Import: i})
		}),
		proto.WithMessage(func(message *proto.Message) {
			ret.Message = append(ret.Message, Message{Message: message})
		}),
		proto.WithPackage(func(p *proto.Package) {
			ret.Package = Package{Package: p}
		}),
		proto.WithService(func(service *proto.Service) {
			serv := Service{Service: service}
			elements := service.Elements
			for _, el := range elements {
				v, _ := el.(*proto.RPC)
				if v == nil {
					continue
				}
				serv.RPC = append(serv.RPC, &RPC{RPC: v})
			}

			serviceList = append(serviceList, serv)
		}),
		proto.WithOption(func(option *proto.Option) {
			if option.Name == "go_package" {
				ret.GoPackage = option.Constant.Source
			}
		}),
	)

	name := filepath.Base(abs)

	if len(ret.GoPackage) == 0 {
		ret.GoPackage = ret.Package.Name
	}
	ret.PbPackage = GoSanitized(filepath.Base(ret.GoPackage))
	ret.Src = abs
	ret.Name = name
	if len(serviceList) > 0 {
		ret.Service = serviceList[0]
	}

	return ret, nil
}

func GoSanitized(s string) string {
	// Sanitize the input to the set of valid characters,
	// which must be '_' or be in the Unicode L or N categories.
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	// Prepend '_' in the event of a Go keyword conflict or if
	// the identifier is invalid (does not start in the Unicode L category).
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}
