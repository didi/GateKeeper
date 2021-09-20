package handler

import (
	"github.com/didi/gatekeeper/public"
	"sync"
	"time"
)

var ServiceCounterHandler *FlowCounter

type FlowCounter struct {
	RedisFlowCountMap map[string]*DistributedCountService
	Locker            sync.RWMutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap: map[string]*DistributedCountService{},
		Locker:            sync.RWMutex{},
	}
}

func init() {
	ServiceCounterHandler = NewFlowCounter()
	ServiceManagerHandler.Regist(ServiceCounterHandler)
}

func (counter *FlowCounter) Update(e *ServiceEvent) {
	for _, service := range e.AddService {
		counter.GetCounter(public.FlowServicePrefix + service.Info.ServiceName)
	}
	for _, item := range counter.RedisFlowCountMap {
		for _, service := range e.DeleteService {
			if item.Name == public.FlowServicePrefix+service.Info.ServiceName {
				item.Close()
				delete(counter.RedisFlowCountMap, item.Name)
			}
		}
	}
}

func (counter *FlowCounter) GetCounter(name string) (*DistributedCountService, error) {
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	if item, ok := counter.RedisFlowCountMap[name]; ok {
		return item, nil
	}
	newCounter := NewDistributedCountService(name, 1*time.Second)
	counter.RedisFlowCountMap[name] = newCounter
	return newCounter, nil
}
