package gentemplate

import (
	_ "embed"
)

//go:embed app.tpl
var AppTpl string

//go:embed ep.go.tpl
var EPTpl string

//go:embed logic.tpl
var LogicTpl string

//go:embed server.tpl
var ServerTpl string

//go:embed message.tpl
var MessageTpl string

//go:embed svc.tpl
var SvcTpl string

//go:embed go.mod.tpl
var GoModTpl string

//go:embed proto.tpl
var ProtoTpl string

//go:embed config.tpl
var ConfigTpl string
