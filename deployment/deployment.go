package deployment

import (
	_ "embed"
	"os"

	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/utils/templatex"
	"github.com/pterm/pterm"
)

//go:embed Deployment.tpl
var DeploymentTpl string

type DeploymentInfo struct {
	Name     string
	Port     int
	Protocol string
	Svc      bool
}

type DeploymentHelper struct {
	Infos     []DeploymentInfo
	Svc       bool // use the svc not hostport
	Config    *config.Config
	Name      string
	Namespace string
}

func NewDeploymentHelper(config *config.Config, svc bool, name string, namespace string) *DeploymentHelper {
	d := &DeploymentHelper{
		Svc:       svc,
		Config:    config,
		Name:      name,
		Namespace: namespace,
	}
	if len(d.Name) == 0 {
		d.Name = "gogs"
	}
	if len(d.Namespace) == 0 {
		d.Namespace = "gogs"
	}
	d.Infos = append(d.Infos, DeploymentInfo{
		Name:     "admin-tcp",
		Port:     config.AdminPort,
		Protocol: "TCP",
		Svc:      svc,
	})

	return d
}

func (d *DeploymentHelper) AddAcceptor(config *acceptor.AcceptroConfig) {
	if config.UdpPort != 0 {
		d.Infos = append(d.Infos, DeploymentInfo{
			Name:     config.Name + "-udp",
			Port:     config.UdpPort,
			Protocol: "UDP",
			Svc:      d.Svc,
		})
	}
	if config.HttpPort != 0 {
		d.Infos = append(d.Infos, DeploymentInfo{
			Name:     config.Name + "-tcp",
			Port:     config.HttpPort,
			Protocol: "TCP",
			Svc:      d.Svc,
		})
	}
}

func (d *DeploymentHelper) Generate() error {
	fileName := "Deployment.yaml"
	if _, err := os.Stat(fileName); err == nil {
		pterm.Error.Printfln(fileName + " already exists")
		return nil
	}
	data := map[string]interface{}{}
	data["Deployments"] = d.Infos
	data["Svc"] = d.Svc
	data["HealthPort"] = d.Config.AdminPort
	data["PodName"] = d.Name
	data["PodNamespace"] = d.Namespace

	if err := templatex.With("gogs").Parse(DeploymentTpl).SaveTo(data, fileName, false); err != nil {
		pterm.Error.Printfln("generate file error :" + err.Error())
		return err
	}

	pterm.Success.Printfln("generate file success " + fileName + " and you should edit the container image with your own image")

	return nil
}
