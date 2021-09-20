package model

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

type ServiceDetail struct {
	Info *ServiceInfo `json:"info" description:"基本信息"`
	PluginConf *simplejson.Json `json:"plugin_conf" description:"plugin_conf"`
}

type UpstreamConfig struct {
	Schema   string
	IpList   []string
	IpWeight map[string]string
}

func GetUpstreamConfigFromString(upstreamList string) (*UpstreamConfig, error) {
	config := &UpstreamConfig{}
	if upstreamList == "" {
		return config, nil
	}
	tmpLine := strings.Split(upstreamList, "\n")
	ipList := []string{}
	ipConf := map[string]string{}
	for _, tmp := range tmpLine {
		r, _ := regexp.Compile("^(.*://)(.*?)\\s(.*?)$")
		submatch := r.FindStringSubmatch(tmp)
		if len(submatch) != 4 {
			return nil, errors.New("upstream_list format error")
		}
		config.Schema = submatch[1]
		ipList = append(ipList, submatch[2])
		ipConf[submatch[2]] = submatch[3]
	}
	config.IpList = ipList
	config.IpWeight = ipConf
	return config, nil
}

func GetServiceDetailFromGinContext(c *gin.Context) (*ServiceDetail, error) {
	serverInterface, ok := c.Get(public.ServiceDetailContextKey)
	if !ok {
		return nil, errors.New("service not found")
	}
	return serverInterface.(*ServiceDetail), nil
}
