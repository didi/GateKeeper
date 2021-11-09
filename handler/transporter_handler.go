package handler

import (
	"github.com/didi/gatekeeper/model"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var TransportorHandler *Transportor

type Transportor struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
	UpdateAt    time.Time
}

func NewTransportor() *Transportor {
	return &Transportor{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

func init() {
	TransportorHandler = NewTransportor()
	ServiceManagerHandler.Regist(TransportorHandler)
}

func (t *Transportor) Update(e *ServiceEvent) {
	for _, service := range e.AddService {
		t.GetTrans(service)
	}
	for _, service := range e.UpdateService {
		t.GetTrans(service)
	}
	newSlice := []*TransportItem{}
	for _, tItem := range t.TransportSlice {
		matched := false
		for _, service := range e.DeleteService {
			if tItem.ServiceName == service.Info.ServiceName {
				matched = true
			}
		}
		if matched {
			delete(t.TransportMap, tItem.ServiceName)
		} else {
			newSlice = append(newSlice, tItem)
		}
	}
	t.TransportSlice = newSlice
}

func (t *Transportor) GetTrans(service *model.ServiceDetail) (*http.Transport, error) {
	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName && transItem.UpdateAt == service.Info.UpdatedAt {
			return transItem.Trans, nil
		}
	}
	idleNumStr := service.PluginConf.GetPath("http_upstream_transport", "http_upstream_connection_idle_num").MustString()
	connectTimeoutStr := service.PluginConf.GetPath("http_upstream_transport", "http_upstream_connection_timeout").MustString()
	headerTimeoutStr := service.PluginConf.GetPath("http_upstream_transport", "http_upstream_header_timeout").MustString()
	idleNum,_:=strconv.ParseInt(idleNumStr,10,64)
	connectTimeout,_:=strconv.ParseInt(connectTimeoutStr,10,64)
	headerTimeout,_:=strconv.ParseInt(headerTimeoutStr,10,64)
	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(connectTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext, //3次握手超时设置
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          int(idleNum),
		WriteBufferSize:       1 << 18, //256m
		ReadBufferSize:        1 << 18, //256m
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(headerTimeout) * time.Second, //请求响应超时
	}

	matched := false
	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName {
			matched = true
			transItem.Trans = trans
			transItem.UpdateAt = service.Info.UpdatedAt
		}
	}
	if !matched {
		transItem := &TransportItem{
			Trans:       trans,
			ServiceName: service.Info.ServiceName,
			UpdateAt:    service.Info.UpdatedAt,
		}
		t.TransportSlice = append(t.TransportSlice, transItem)
		t.Locker.Lock()
		defer t.Locker.Unlock()
		t.TransportMap[service.Info.ServiceName] = transItem
	}
	return trans, nil
}
