package http_proxy_router

import (
	"context"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/http_proxy_middleware"
	"net/http"
	"strings"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	r := InitRouter(http_proxy_middleware.HTTPRecoveryMiddleware(), http_proxy_middleware.HTTPRequestLogger())
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	httpAddr := lib.GetStringConf("proxy.http.addr")
	if strings.HasPrefix(httpAddr, ":") {
		httpAddr = "http://127.0.0.1" + httpAddr
	}
	log.Info().Msg(lib.Purple(fmt.Sprintf("start HTTP proxy service [%s]\n", httpAddr)))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Msg(lib.Purple(fmt.Sprintf("failed to start HTTPS proxy service [%s] %v", httpAddr, err)))
	}
}

func HttpsServerRun() {
	r := InitRouter(http_proxy_middleware.HTTPRecoveryMiddleware(), http_proxy_middleware.HTTPRequestLogger())
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	httpsAddr := lib.GetStringConf("proxy.https.addr")
	if strings.HasPrefix(httpsAddr, ":") {
		httpsAddr = "https://127.0.0.1" + httpsAddr
	}
	log.Info().Msg(lib.Purple(fmt.Sprintf("start HTTPS proxy service [%s]", httpsAddr)))
	if err := HttpsSrvHandler.ListenAndServeTLS("./cert_file/server.crt", "./cert_file/server.key"); err != nil && err != http.ErrServerClosed {
		log.Error().Msg(lib.Purple(fmt.Sprintf("failed to start HTTPS proxy service [%s] %v", httpsAddr, err)))
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Error().Msg(lib.Purple(fmt.Sprintf("http_proxy_stop err:%v", err)))
	}
	log.Error().Msg(lib.Purple(fmt.Sprintf("http_proxy_stop %v stopped", lib.GetStringConf("proxy.http.addr"))))
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Error().Msg(lib.Purple(fmt.Sprintf("https_proxy_stop err:%v", err)))
	}
	log.Error().Msg(lib.Purple(fmt.Sprintf("https_proxy_stop %v stopped", lib.GetStringConf("proxy.http.addr"))))
}
