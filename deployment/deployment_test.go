package deployment

import (
	_ "embed"
	"os"
	"testing"

	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/utils/filex"
)

func TestDeploymentHelper_Generate(t *testing.T) {
	type fields struct {
		Infos     []DeploymentInfo
		Svc       bool
		Config    *config.Config
		Name      string
		Namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test1.yaml",
			fields: fields{
				Infos: []DeploymentInfo{
					{
						Name:     "admin-tcp",
						Port:     8080,
						Protocol: "TCP",
						Svc:      true,
					},
					{
						Name:     "admin-udp",
						Port:     8081,
						Protocol: "UDP",
						Svc:      true,
					},
				},
				Config: &config.Config{
					AdminPort: 8080,
				},
				Svc:       true,
				Name:      "testname",
				Namespace: "testnamespace",
			},
			wantErr: false,
		},
		{
			name: "test2.yaml",
			fields: fields{
				Infos: []DeploymentInfo{
					{
						Name:     "admin-tcp",
						Port:     8080,
						Protocol: "TCP",
						Svc:      false,
					},
					{
						Name:     "admin-udp",
						Port:     8081,
						Protocol: "UDP",
						Svc:      false,
					},
				},
				Config: &config.Config{
					AdminPort: 8080,
				},
				Svc:       false,
				Name:      "testname",
				Namespace: "testnamespace",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeploymentHelper{
				Infos:     tt.fields.Infos,
				Svc:       tt.fields.Svc,
				Config:    tt.fields.Config,
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
			}
			if err := d.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("DeploymentHelper.Generate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if ok := filex.IsFileEqual("testdata/"+tt.name, "Deployment.yaml"); !ok {
				t.Errorf("Init.Generate() error = %s is not equal to %s", "testdata/"+tt.name, "Deployment.yaml")
			}
			_ = os.Remove("deployment.yaml")
		})
	}
}
