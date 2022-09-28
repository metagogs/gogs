package gentemplate

import (
	_ "embed"
)

//go:embed docker.tpl
var DockerTpl string
