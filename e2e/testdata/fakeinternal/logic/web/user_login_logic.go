package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
)

type UserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	UID string `json:"uid"`
}

func (l *UserLoginLogic) Handler(c *gin.Context) {
	var msg UserLoginRequest
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	if len(msg.Username) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "username is empty",
		})
		return
	}

	response := UserLoginResponse{
		UID: l.svcCtx.SF.Generate().String(),
	}

	l.svcCtx.PlayerManagaer.CreateUser(response.UID, msg.Username)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})

}
