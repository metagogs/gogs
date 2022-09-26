package gentemplate

import (
	_ "embed"
)

//go:embed Gogs/Codec.tpl
var CodecTpl string

//go:embed Gogs/Common.tpl
var CommonTpl string

//go:embed Gogs/EventsManager.tpl
var EventsManagerTpl string

//go:embed Gogs/ICodec.tpl
var ICodecTpl string

//go:embed Gogs/Messages.tpl
var MessagesTpl string

//go:embed Gogs/Packet.tpl
var PacketTpl string

//go:embed Register.tpl
var RegisterTpl string
