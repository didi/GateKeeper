package http_proxy_router

import (
	"github.com/didi/gatekeeper/dashboard_controller"
	"github.com/didi/gatekeeper/http_proxy_middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	oauth := router.Group("/oauth")
	oauth.Use(http_proxy_middleware.TranslationMiddleware())
	{
		dashboard_controller.OAuthRegister(oauth)
	}
	router.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),
		http_proxy_middleware.HTTPFlowCountMiddleware(),
		http_proxy_middleware.HTTPFlowLimitMiddleware(),
		http_proxy_middleware.HTTPWhiteListMiddleware(),
		http_proxy_middleware.HTTPBlackListMiddleware(),
		http_proxy_middleware.HTTPJwtAuthTokenMiddleware(),
		http_proxy_middleware.HTTPJwtFlowCountMiddleware(),
		http_proxy_middleware.HTTPJwtFlowLimitMiddleware(),
		http_proxy_middleware.HTTPHeaderTransferMiddleware(),
		http_proxy_middleware.HTTPStripUriMiddleware(),
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),
		http_proxy_middleware.HTTPReverseProxyMiddleware())
	return router
}
