package middleware

import (
	"context"
	"fmt"
	"github.com/didichuxing/gatekeeper/dao"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/didichuxing/gatekeeper/tcpproxy"
	"net"
	"os"
	"strings"
)

//TCPLimit 必须返回一个实现该方法的回调函数
func TCPLimit(module *dao.GatewayModule) tcpproxy.HandlerFunc {
	moduleName := module.Base.Name
	qps := module.AccessControl.ClientFlowLimit
	return func(next tcpproxy.Target) tcpproxy.Target {
		return tcpproxy.TCPHandlerFunc(func(c net.Conn) {
			//统计
			counter := public.FlowCounterHandler.GetRequestCounter(moduleName)
			counter.Increase(context.Background(), c.RemoteAddr().String())
			remoteIP := public.Substr(c.RemoteAddr().String(), 0,
				int64(strings.Index(c.RemoteAddr().String(), ":")))
			if qps > 0 {
				//限流
				limiter := public.FlowLimiterHandler.GetModuleIPVisitor(moduleName+"_"+remoteIP, qps)
				if limiter.Allow() == false {
					errmsg := fmt.Sprintf("moduleName:%s remoteIP：%s, QPS limit : %d, %d", moduleName, remoteIP, int64(limiter.Limit()), limiter.Burst())
					public.ContextWarning(context.Background(), service.DLTagAccessControlFailure, map[string]interface{}{
						"msg":        errmsg,
						"ip":         remoteIP,
						"moduleName": module.Base.Name,
					})
					c.Close()
					return
				}
			}

			//白名单限制
			blackList := strings.Split(module.AccessControl.BlackList, ",")
			whiteList := strings.Split(module.AccessControl.WhiteList, ",")
			whiteHostname := strings.Split(module.AccessControl.WhiteHostName, ",")
			hostname, _ := os.Hostname()
			authed := false
			if module.AccessControl.Open == 0 {
				authed = true
				public.ContextNotice(context.Background(), service.DLTagAccessControlUndef, map[string]interface{}{
					"msg":      "access_control_close",
					"clientip": remoteIP,
				})
			} else if public.AuthIPList(remoteIP, blackList) {
				errmsg := "client ip in blacklist"
				public.ContextWarning(context.Background(), service.DLTagAccessControlFailure, map[string]interface{}{
					"msg":       errmsg,
					"ip":        remoteIP,
					"blackList": blackList,
				})
			} else if public.AuthIPList(remoteIP, whiteList) {
				authed = true
				public.ContextNotice(context.Background(), service.DLTagAccessControlUndef, map[string]interface{}{
					"msg":       "clientip_in_whitelist",
					"clientip":  remoteIP,
					"whitelist": whiteList,
				})
			} else if public.InStringList(hostname, whiteHostname) {
				authed = true
				public.ContextNotice(context.Background(), service.DLTagAccessControlUndef, map[string]interface{}{
					"msg":           "hostname_in_whitehostlist",
					"hostname":      hostname,
					"whitehostlist": whiteHostname,
				})
			}
			if authed == false {
				c.Close()
				return
			}
			next.HandleConn(c)
		})
	}
}
