package dashboard_router

import (
	"github.com/didi/gatekeeper/dashboard_controller"
	"github.com/didi/gatekeeper/dashboard_middleware"
	"github.com/didi/gatekeeper/docs"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adminLoginRouter := router.Group("/admin_login")
	store, err := sessions.NewRedisStore(10, "tcp", lib.GetStringConf("base.session.redis_server"), lib.GetStringConf("base.session.redis_password"), []byte("secret"))
	if err != nil {
		log.Fatalf("sessions.NewRedisStore err:%v", err)
	}
	adminLoginRouter.Use(
		sessions.Sessions("mysession", store),
		dashboard_middleware.RecoveryMiddleware(),
		dashboard_middleware.RequestLogger("gatekeeper_dashboard"),
		dashboard_middleware.TranslationMiddleware())
	{
		dashboard_controller.AdminLoginRegister(adminLoginRouter)
	}

	adminRouter := router.Group("/admin")
	adminRouter.Use(
		sessions.Sessions("mysession", store),
		dashboard_middleware.RecoveryMiddleware(),
		dashboard_middleware.RequestLogger("gatekeeper_dashboard"),
		dashboard_middleware.SessionAuthMiddleware(),
		dashboard_middleware.TranslationMiddleware())
	{
		dashboard_controller.AdminRegister(adminRouter)
	}

	serviceRouter := router.Group("/service")
	serviceRouter.Use(
		sessions.Sessions("mysession", store),
		dashboard_middleware.RecoveryMiddleware(),
		dashboard_middleware.RequestLogger("gatekeeper_dashboard"),
		dashboard_middleware.SessionAuthMiddleware(),
		dashboard_middleware.TranslationMiddleware())
	{
		dashboard_controller.ServiceRegister(serviceRouter)
	}

	appRouter := router.Group("/app")
	appRouter.Use(
		sessions.Sessions("mysession", store),
		dashboard_middleware.RecoveryMiddleware(),
		dashboard_middleware.RequestLogger("gatekeeper_dashboard"),
		dashboard_middleware.SessionAuthMiddleware(),
		dashboard_middleware.TranslationMiddleware())
	{
		dashboard_controller.APPRegister(appRouter)
	}

	dashRouter := router.Group("/dashboard")
	dashRouter.Use(
		sessions.Sessions("mysession", store),
		dashboard_middleware.RecoveryMiddleware(),
		dashboard_middleware.RequestLogger("gatekeeper_dashboard"),
		dashboard_middleware.SessionAuthMiddleware(),
		dashboard_middleware.TranslationMiddleware())
	{
		dashboard_controller.DashboardRegister(dashRouter)
	}

	router.Static("/dist", "./dist")
	return router
}
