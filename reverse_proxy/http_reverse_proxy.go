package reverse_proxy

import (
	"bytes"
	"github.com/didi/gatekeeper/dashboard_middleware"
	"github.com/didi/gatekeeper/load_balance"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewLoadBalanceReverseProxy(c *gin.Context, lb *load_balance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil || nextAddr == "" {
			panic("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			buffer := bytes.NewBufferString(targetQuery)
			buffer.WriteString(req.URL.RawQuery)
			req.URL.RawQuery = buffer.String()
		} else {
			buffer := bytes.NewBufferString(targetQuery)
			buffer.WriteString("&")
			buffer.WriteString(req.URL.RawQuery)
			req.URL.RawQuery = buffer.String()
		}
	}
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		dashboard_middleware.ResponseError(c, 999, err)
	}
	return &httputil.ReverseProxy{Director: director, Transport: trans, ErrorHandler: errFunc}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		buffer := bytes.NewBufferString(a)
		buffer.WriteString(b[1:])
		return buffer.String()
	case !aslash && !bslash:
		buffer := bytes.NewBufferString(a)
		buffer.WriteString("/")
		buffer.WriteString(b)
		return buffer.String()
	}
	buffer := bytes.NewBufferString(a)
	buffer.WriteString(b)
	return buffer.String()
}
