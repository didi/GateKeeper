package dashboard_router

import (
	"context"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.GetStringConf("base.base.debug_mode"))
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	go func() {
		log.Info().Msg(lib.Purple(fmt.Sprintf("start HTTP control service [http://127.0.0.1%s/dist/]", lib.GetStringConf("base.http.addr"))))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Error().Msg(lib.Purple(fmt.Sprintf("failed to start HTTP service service [%s] %v", lib.GetStringConf("base.http.addr"), err)))
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Error().Msg(lib.Purple(fmt.Sprintf("HttpServerStop err:%v", err)))
	}
	log.Error().Msg(lib.Purple(fmt.Sprintf("stop HTTP control service [%s]", lib.GetStringConf("base.http.addr"))))
}
