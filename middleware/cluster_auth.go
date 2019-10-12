package middleware

import (
	"github.com/didi/gatekeeper/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

//ClusterAuth 集群验证中间件
func ClusterAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterList := lib.GetStringConf("base.cluster.cluster_list")
		matchFlag := false
		ipList := strings.Split(clusterList, ",")
		ipList = append(ipList, "127.0.0.1")
		for _, host := range ipList {
			if c.ClientIP() == host {
				matchFlag = true
			}
		}
		if !matchFlag {
			public.ResponseError(c, http.StatusBadRequest, errors.New("ClusterAuth error"))
			return
		}
		c.Next()
	}
}
