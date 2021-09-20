package handler

import (
	"github.com/didi/gatekeeper/public"
	"log"
	"sync"
)

var AppCounterHandler *FlowAppCounter

type FlowAppCounter struct {
	FlowCounter
}

func NewFlowAppCounter() *FlowAppCounter {
	return &FlowAppCounter{
		FlowCounter{
			RedisFlowCountMap: map[string]*DistributedCountService{},
			Locker:            sync.RWMutex{},
		},
	}
}

func (counter *FlowAppCounter) Update(e *AppEvent) {
	log.Printf("FlowAppCounter.Update\n")
	for _, app := range e.AddApp {
		counter.GetCounter(public.FlowAppPrefix + app.AppID)
	}
	for _, item := range counter.RedisFlowCountMap {
		for _, app := range e.DeleteApp {
			if item.Name == public.FlowAppPrefix+app.AppID {
				item.Close()
				delete(counter.RedisFlowCountMap, item.Name)
			}
		}
	}
}

func init() {
	AppCounterHandler = NewFlowAppCounter()
}
