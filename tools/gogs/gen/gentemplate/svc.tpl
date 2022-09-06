package svc

import (
	"github.com/metagogs/gogs"
)

type ServiceContext struct {
	*gogs.App
}

func NewServiceContext(app *gogs.App) *ServiceContext {
	return &ServiceContext{
		App: app,
	}
}
