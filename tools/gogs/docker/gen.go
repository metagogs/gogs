package docker

import (
	"github.com/metagogs/gogs/tools/gogs/docker/gentemplate"
	"github.com/metagogs/gogs/utils/templatex"
	"github.com/pterm/pterm"
)

type DockerGen struct {
	cnproxy bool
}

func NewDockerGen(cnproxy bool) (*DockerGen, error) {
	return &DockerGen{
		cnproxy: cnproxy,
	}, nil
}

func (g *DockerGen) Generate() error {

	if err := g.docker(); err != nil {
		return err
	}

	return nil
}

func (g *DockerGen) docker() error {
	data := map[string]interface{}{}
	data["Proxy"] = g.cnproxy

	if err := templatex.With("gogs").Parse(gentemplate.DockerTpl).SaveTo(data, "Dockerfile", true); err != nil {
		pterm.Error.Printfln("generate file error Dockerfile:" + err.Error())
		return err
	}

	return nil
}
