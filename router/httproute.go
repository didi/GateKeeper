package router

import (
	"github.com/didichuxing/gatekeeper/controller"
	"github.com/didichuxing/gatekeeper/middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

//InitRouter 声明http路由
func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Recovery())

	//admin
	admin := router.Group("/admin")
	admin.Use(middleware.RequestTraceLog())
	{
		controller.AdminRegister(admin)
	}

	//assets
	router.Static("/assets", "./tmpl/green/assets")

	//gateway
	gateway := controller.Gateway{}
	router.GET("/ping", gateway.Ping)

	//cluster
	csr:=router.Group("/")
	csr.Use(middleware.ClusterAuth())
	csr.GET("/reload", gateway.Reload)

	gw:=router.Group(lib.GetStringConf("base.http.route_prefix"))
	gw.Use(
		middleware.RequestTraceLog(),
		middleware.MatchRule(),
		middleware.AccessControl(),
		middleware.HTTPLimit(),
		//todo 拓展中间件
		middleware.LoadBalance())
	{
		gw.GET("/*action", gateway.Index)
		gw.POST("/*action", gateway.Index)
		gw.DELETE("/*action", gateway.Index)
		gw.OPTIONS("/*action", gateway.Index)
	}
	return router
}