package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
)

var ServiceManagerHandler *ServiceManager = NewServiceManager()

func ServiceManagerInit() {
	ServiceManagerHandler = NewServiceManager()
}

//通知事件
type ServiceEvent struct {
	DeleteService []*model.ServiceDetail
	AddService    []*model.ServiceDetail
	UpdateService []*model.ServiceDetail
}

//观察者接口
type ServiceObserver interface {
	Update(*ServiceEvent)
}

//被观察者接口
type ServiceSubject interface {
	Regist(ServiceObserver)
	Deregist(ServiceObserver)
	Notify(*ServiceEvent)
}

func (s *ServiceManager) Regist(ob ServiceObserver) {
	s.Lock()
	defer s.Unlock()
	s.Observers[ob] = true
}

func (s *ServiceManager) Deregist(ob ServiceObserver) {
	s.Lock()
	defer s.Unlock()
	delete(s.Observers, ob)
}

func (s *ServiceManager) Notify(e *ServiceEvent) {
	s.RLock()
	defer s.RUnlock()
	for ob, _ := range s.Observers {
		ob.Update(e)
	}
}

type ServiceManager struct {
	ServiceMap   map[string]*model.ServiceDetail
	ServiceSlice []*model.ServiceDetail
	err          error
	UpdateAt     time.Time
	Observers    map[ServiceObserver]bool
	sync.RWMutex
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*model.ServiceDetail{},
		ServiceSlice: []*model.ServiceDetail{},
		Observers:    map[ServiceObserver]bool{},
	}
}

func (s *ServiceManager) GetTcpServiceList() []*model.ServiceDetail {
	list := []*model.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.ServiceType == public.LoadTypeTCP {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *ServiceManager) GetGrpcServiceList() []*model.ServiceDetail {
	list := []*model.ServiceDetail{}
	for _, serverItem := range s.ServiceSlice {
		tempItem := serverItem
		if tempItem.Info.ServiceType == public.LoadTypeGRPC {
			list = append(list, tempItem)
		}
	}
	return list
}

func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*model.ServiceDetail, error) {
	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.ServiceType != public.LoadTypeHTTP {
			continue
		}
		hosts := strings.Split(serviceItem.Info.HTTPHosts, "\n")
		paths := strings.Split(serviceItem.Info.HTTPPaths, "\n")
		for _, path := range paths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				//log.Info().Msgf("new_matched, path=%s", paths)
				return serviceItem, nil
			}
		}
		if serviceItem.Info.HTTPHosts != "" {
			for _, path := range paths {
				if strings.HasPrefix(c.Request.URL.Path, path) && public.InArrayString(c.Request.Host, hosts) {
					return serviceItem, nil
				}
			}
		}
	}
	return nil, errors.New("not matched service")
}

func (s *ServiceManager) LoadService() *ServiceManager {
	ns := NewServiceManager()
	defer func() {
		if ns.err != nil {
			log.Error().Msgf("load service config error:%v", ns.err)
		}
	}()
	serviceInfo := &model.ServiceInfo{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	tx, err := lib.GetGormPool("default")
	if err != nil {
		ns.err = err
		return ns
	}
	params := &model.ServiceListInput{PageNo: 1, PageSize: 99999}
	list, _, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		ns.err = err
		return ns
	}
	for _, listItem := range list {
		tmpItem := listItem
		serviceDetail, err := tmpItem.ServiceDetail(c, tx, &tmpItem)
		if err != nil {
			ns.err = err
			return ns
		}
		ns.ServiceMap[listItem.ServiceName] = serviceDetail
		ns.ServiceSlice = append(ns.ServiceSlice, serviceDetail)
		if listItem.UpdatedAt.Unix() > ns.UpdateAt.Unix() {
			ns.UpdateAt = listItem.UpdatedAt
		}
	}
	return ns
}

func (s *ServiceManager) Load() error {
	ns := s.LoadService()
	if ns.err != nil {
		return ns.err
	}
	s.ServiceSlice = ns.ServiceSlice
	s.ServiceMap = ns.ServiceMap
	s.UpdateAt = ns.UpdateAt
	e := &ServiceEvent{AddService: ns.ServiceSlice}
	s.Notify(e)
	return s.err
}

func (s *ServiceManager) LoadAndWatch() error {
	log.Info().Msg(lib.Purple("watching load service config from resource"))
	ns := s.LoadService()
	if ns.err != nil {
		return ns.err
	}
	s.ServiceSlice = ns.ServiceSlice
	s.ServiceMap = ns.ServiceMap
	s.UpdateAt = ns.UpdateAt
	e := &ServiceEvent{AddService: ns.ServiceSlice}
	s.Notify(e)
	go func() {
		for true {
			time.Sleep(10 * time.Second)
			ns := s.LoadService()
			if ns.err != nil {
				log.Info().Msg(lib.Purple(fmt.Sprintf("load service err:%v update in:%v", ns.err)))
				continue
			}
			if ns.UpdateAt != s.UpdateAt || len(ns.ServiceSlice) != len(s.ServiceSlice) {
				e := &ServiceEvent{}
				for _, service := range s.ServiceSlice {
					matched := false
					for _, newService := range ns.ServiceSlice {
						if service.Info.ServiceName == newService.Info.ServiceName {
							matched = true
						}
					}
					if !matched {
						e.DeleteService = append(e.DeleteService, service)
					}
				}
				for _, newService := range ns.ServiceSlice {
					matched := false
					for _, service := range s.ServiceSlice {
						if service.Info.ServiceName == newService.Info.ServiceName {
							matched = true
						}
					}
					if !matched {
						e.AddService = append(e.AddService, newService)
					}
				}
				for _, newService := range ns.ServiceSlice {
					matched := false
					for _, service := range s.ServiceSlice {
						if service.Info.ServiceName == newService.Info.ServiceName && service.Info.UpdatedAt != newService.Info.UpdatedAt {
							matched = true
						}
					}
					if matched {
						e.UpdateService = append(e.UpdateService, newService)
					}
				}
				for _, item := range e.DeleteService {
					log.Info().Msg(lib.Purple(fmt.Sprintf("found config delete service[%v] update_time[%v]", item.Info.ServiceName, ns.UpdateAt.Format(lib.TimeFormat))))
				}
				for _, item := range e.AddService {
					log.Info().Msg(lib.Purple(fmt.Sprintf("found config add service[%v] update_time[%v]", item.Info.ServiceName, ns.UpdateAt.Format(lib.TimeFormat))))
				}
				for _, item := range e.UpdateService {
					log.Info().Msg(lib.Purple(fmt.Sprintf("found config update service[%v] update_time[%v]", item.Info.ServiceName, ns.UpdateAt.Format(lib.TimeFormat))))
				}
				s.ServiceSlice = ns.ServiceSlice
				s.ServiceMap = ns.ServiceMap
				s.UpdateAt = ns.UpdateAt
				//log.Info().Msg(lib.Purple(fmt.Sprintf("e:%v", e)))
				s.Notify(e)
			}
		}
	}()
	return s.err
}
