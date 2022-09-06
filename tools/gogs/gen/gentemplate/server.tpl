package server

import (
	"context"

	"{{.BasePackage}}/{{.Package}}"
	{{range .Components}} "{{.BasePackage}}/internal/logic/{{.Name | ToLower}}" 
	{{end}}
	"{{.BasePackage}}/internal/svc"
	"github.com/metagogs/gogs/session"
)

type Server struct {
	svcCtx *svc.ServiceContext
}

func NewServer(svcCtx *svc.ServiceContext) *Server {
	return &Server{
		svcCtx: svcCtx,
	}
}



{{range .Components}}
{{range .Fields}}

{{if not .ServerMessage}}
func (gogs *Server) {{.Name}}(ctx context.Context, s *session.Session, in *{{.Package}}.{{.Name}}) {
	l := {{ .ComponentName | ToLower  }}.New{{.Name}}Logic(ctx, gogs.svcCtx, s)
	l.Handler(in)
}
{{end}}

{{end}}
{{end}}


