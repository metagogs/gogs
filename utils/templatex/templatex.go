package templatex

import (
	"bytes"
	goformat "go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	nameUtil "github.com/metagogs/gogs/utils/name"
)

const regularPerm = 0o666

// DefaultTemplate is a tool to provides the text/template operations
type DefaultTemplate struct {
	name     string
	text     string
	goFmt    bool
	funcMap  template.FuncMap
	savePath string
}

// With returns a instance of DefaultTemplate
func With(name string) *DefaultTemplate {
	d := &DefaultTemplate{
		name: name,
	}
	d.funcMap = template.FuncMap{
		"ToUpper":   strings.ToUpper,
		"ToLower":   strings.ToLower,
		"CamelCase": nameUtil.CamelCase,
	}
	return d
}

// Parse accepts a source template and returns DefaultTemplate
func (t *DefaultTemplate) Parse(text string) *DefaultTemplate {
	t.text = text
	return t
}

func (t *DefaultTemplate) Funcs(f template.FuncMap) *DefaultTemplate {
	t.funcMap = f
	return t
}

// GoFmt sets the value to goFmt and marks the generated codes will be formatted or not
func (t *DefaultTemplate) GoFmt(format bool) *DefaultTemplate {
	t.goFmt = format
	return t
}

// SaveTo writes the codes to the target path
func (t *DefaultTemplate) SaveTo(data interface{}, path string, forceUpdate bool) error {
	if FileExists(path) && !forceUpdate {
		return nil
	}

	output, err := t.Execute(data)
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)

	return os.WriteFile(path, output.Bytes(), regularPerm)
}

// Execute returns the codes after the template executed
func (t *DefaultTemplate) Execute(data interface{}) (*bytes.Buffer, error) {
	tem, err := template.New(t.name).Funcs(t.funcMap).Parse(t.text)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = tem.Execute(buf, data); err != nil {
		return nil, err
	}

	if !t.goFmt {
		return buf, nil
	}

	formatOutput, err := goformat.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	buf.Reset()
	buf.Write(formatOutput)
	return buf, nil
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
