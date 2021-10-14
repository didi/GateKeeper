package handler

import (
	"github.com/didi/gatekeeper/load_balance"
	"github.com/didi/gatekeeper/model"
	"sync"
	"time"
)

var LoadBalancerHandler *LoadBalancer

type LoadBalancer struct {
	LoadBanlanceMap   map[string]*LoadBalancerItem
	LoadBanlanceSlice []*LoadBalancerItem
	Locker            sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBanlance *load_balance.LoadBalance
	ServiceName  string
	UpdatedAt    time.Time
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBanlanceMap:   map[string]*LoadBalancerItem{},
		LoadBanlanceSlice: []*LoadBalancerItem{},
		Locker:            sync.RWMutex{},
	}
}

func init() {
	LoadBalancerHandler = NewLoadBalancer()
	ServiceManagerHandler.Regist(LoadBalancerHandler)
}

func (lbr *LoadBalancer) Update(e *ServiceEvent) {
	for _, service := range e.AddService {
		lbr.GetLoadBalancer(service)
	}
	for _, service := range e.UpdateService {
		lbr.GetLoadBalancer(service)
	}
	newLBSlice := []*LoadBalancerItem{}
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		matched := false
		for _, service := range e.DeleteService {
			if lbrItem.ServiceName == service.Info.ServiceName {
				lbrItem.LoadBanlance.Close()
				matched = true
			}
		}
		if matched {
			delete(lbr.LoadBanlanceMap, lbrItem.ServiceName)
		} else {
			newLBSlice = append(newLBSlice, lbrItem)
		}
	}
	lbr.LoadBanlanceSlice = newLBSlice
}

func (lbr *LoadBalancer) GetLoadBalancer(service *model.ServiceDetail) (*load_balance.LoadBalance, error) {
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName && lbrItem.UpdatedAt == service.Info.UpdatedAt {
			return lbrItem.LoadBanlance, nil
		}
	}
	//fmt.Println("service.Info.LoadBalanceType", service.Info.LoadBalanceType)
	confHandler := load_balance.GetCheckConfigHandler(service.Info.LoadBalanceType)
	checkConf, err := confHandler(service)
	if err != nil {
		return nil, err
	}
	//fmt.Println("service.Info.LoadBalanceStrategy", service.Info.LoadBalanceStrategy)
	lb := load_balance.LoadBanlanceFactorWithStrategy(load_balance.GetLoadBalanceStrategy(service.Info.LoadBalanceStrategy), checkConf)
	matched := false
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName {
			matched = true
			lbrItem.LoadBanlance.Close()
			lbrItem.LoadBanlance = lb
			lbrItem.UpdatedAt = service.Info.UpdatedAt
		}
	}
	if !matched {
		lbItem := &LoadBalancerItem{
			LoadBanlance: lb,
			ServiceName:  service.Info.ServiceName,
			UpdatedAt:    service.Info.UpdatedAt,
		}
		lbr.LoadBanlanceSlice = append(lbr.LoadBanlanceSlice, lbItem)
		lbr.Locker.Lock()
		defer lbr.Locker.Unlock()
		lbr.LoadBanlanceMap[service.Info.ServiceName] = lbItem
	}
	return lb, nil
}
