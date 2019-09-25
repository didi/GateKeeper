package controller

import (
	"github.com/didi/gatekeeper/public"
	"github.com/didi/gatekeeper/service"
	"github.com/gin-gonic/gin"
)

//Gateway struct
type Gateway struct {
}

//Index /index
func (g *Gateway) Index(c *gin.Context) {
	public.ResponseSuccess(c, string("gateway index"))
	return
}

//Ping /ping
func (g *Gateway) Ping(c *gin.Context) {
	public.ResponseSuccess(c, string("gateway pong"))
	return
}

//Reload /reload
func (g *Gateway) Reload(c *gin.Context) {
	service.SysConfMgr.ReloadConfig()
	public.ResponseSuccess(c, string("gateway config loaded"))
	return
}
