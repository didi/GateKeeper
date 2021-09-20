package tcp_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/public"
	"strings"
)

//匹配接入方式 基于请求信息
func TCPBlackListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serviceDetail, err := c.GetServiceDetail()
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		whiteListStr := serviceDetail.PluginConf.GetPath("tcp_whiteblacklist", "ip_white_list").MustString()
		blackListStr := serviceDetail.PluginConf.GetPath("tcp_whiteblacklist", "ip_black_list").MustString()
		if blackListStr == "" {
			c.Next()
			return
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		if whiteListStr == "" && blackListStr != "" {
			if public.InIPSliceStr(clientIP, blackListStr) {
				c.conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
