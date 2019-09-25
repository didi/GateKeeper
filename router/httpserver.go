package router

import (
	"context"
	"github.com/didi/gatekeeper/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	//HTTPSrvHandler 全局变量
	HTTPSrvHandler *http.Server
)

//HTTPServerRun 服务启动
func HTTPServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter()
	HTTPSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				public.SysLogger.Error("HttpServerRun_recover:%v", err)
			}
		}()
		log.Printf(" [INFO] HttpServer %s listening\n",lib.GetStringConf("base.http.addr"))
		if err := HTTPSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" [ERROR] HttpServer %s err:%v\n", lib.GetStringConf("base.http.addr"), err)
		}
	}()
}

//HTTPServerStop 服务停止
func HTTPServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HTTPSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServer err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServer stopped\n")
}
