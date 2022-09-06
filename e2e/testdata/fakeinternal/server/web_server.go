package server

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/logic/web"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
)

type WebServer struct {
	svcCtx *svc.ServiceContext
	engine *gin.Engine
}

func NewWebServer(svcCtx *svc.ServiceContext) *WebServer {
	return &WebServer{
		svcCtx: svcCtx,
	}
}

func (s *WebServer) RegisterHandler(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))
	engine.POST("/user/login", s.UserLogin)
}

func (s *WebServer) UserLogin(c *gin.Context) {
	web.NewUserLoginLogic(context.Background(), s.svcCtx).Handler(c)
}
