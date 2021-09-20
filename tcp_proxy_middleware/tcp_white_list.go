package tcp_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/public"
	"strings"
)

//匹配接入方式 基于请求信息
func TCPWhiteListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serviceDetail, err := c.GetServiceDetail()
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		whiteListStr := serviceDetail.PluginConf.GetPath("tcp_whiteblacklist", "ip_white_list").MustString()
		if whiteListStr != "" {
			if !public.InIPSliceStr(clientIP, whiteListStr) {
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIP)))
				c.Abort()
				return
			}
		}
	
		c.Next()
	}
}
