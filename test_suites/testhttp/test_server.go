package testhttp

import (
	"flag"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"net/http"
	"sync"
	"time"
)

//HTTPDestServer 目标服务管理器
type HTTPDestServer struct {
	mapHTTPDestSvr map[string]*http.Server
	lock           *sync.Mutex
}

//NewTestHTTPDestServer 创建目标服务管理器
func NewTestHTTPDestServer() *HTTPDestServer {
	return &HTTPDestServer{
		mapHTTPDestSvr: map[string]*http.Server{},
		lock:           &sync.Mutex{},
	}
}

//Run 启动服务
func (t *HTTPDestServer) Run(addr string, showLog ...bool) {
	flag.Parse()
	if addr == "" {
		log.Info().Msg(lib.Purple("need addr like :8007"))
		log.Fatal()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", t.getPath)
	mux.HandleFunc("/get_path", t.getPath)
	mux.HandleFunc("/get_host", t.getHost(fmt.Sprintf("127.0.0.1%s", addr)))
	mux.HandleFunc("/ping", t.ping)
	mux.HandleFunc("/goods_list", t.goodsList)
	if len(showLog) == 0 || showLog[0] == true {
		log.Info().Msg(lib.Purple("start test httpserver:" + addr))
	}
	go func() {
		server := &http.Server{
			Addr:         addr,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			Handler:      mux,
		}
		t.lock.Lock()
		t.mapHTTPDestSvr[addr] = server
		t.lock.Unlock()
		err := server.ListenAndServe()
		if err != nil {
			//log.Printf("\nRunHttpDestServer addr:%v -err:%v", addr, err)
		}
	}()
}

//Stop 关闭服务
func (t *HTTPDestServer) Stop(addr string) {
	if server, ok := t.mapHTTPDestSvr[addr]; ok {
		server.Close()
	}
}

func (t *HTTPDestServer) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func (t *HTTPDestServer) getHost(addr string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, addr)
	}
}

func (t *HTTPDestServer) getPath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}

func (t *HTTPDestServer) goodsList(w http.ResponseWriter, r *http.Request) {
	str := "{\"errno\":0,\"errmsg\":\"\",\"data\":{\"list\":[{\"pid\":\"2018103018_9000042581625\",\"name\":\"商品1\",\"city_id\":\"12\"},{\"pid\":\"2018103018_9000042581625\",\"name\":\"商品2\",\"city_id\":\"1\"},{\"pid\":\"2018103018_9000042581625\",\"name\":\"商品3\",\"city_id\":\"13\"},{\"pid\":\"2018103018_9000042581625\",\"name\":\"商品4\",\"city_id\":\"2\"}]}}"
	fmt.Fprintf(w, str)
}
