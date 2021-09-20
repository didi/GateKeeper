package handler

import (
	"fmt"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

var FlowLimiterHandler *FlowLimiter

type FlowLimiter struct {
	FlowLimiterMap map[string]*FlowLimiterItem
	Locker         sync.RWMutex
}

type Limiter interface {
	Allow() bool
}

type FlowLimiterItem struct {
	LastTime      time.Time
	Name          string
	LType         int  //限流类型 0=qps 1=qpm 2=qph
	IsDistributed bool //是否分布式
	Limter        Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap: map[string]*FlowLimiterItem{},
		Locker:         sync.RWMutex{},
	}
}

func CleanupLimiters() {
	for {
		time.Sleep(10 * time.Second)
		for hashName, v := range FlowLimiterHandler.FlowLimiterMap {
			if time.Now().Sub(v.LastTime) > time.Hour {
				delete(FlowLimiterHandler.FlowLimiterMap, hashName)
			}
		}
	}
}

func init() {
	FlowLimiterHandler = NewFlowLimiter()
	go CleanupLimiters()
}

//一个服务,可能设置多种限流策略
func (counter *FlowLimiter) GetLimiter(serverName string, val float64, ltype int, isDistributed bool) (Limiter, error) {
	hashName := fmt.Sprintf("%s_%d_%d_%v", serverName, ltype, val, isDistributed)
	counter.Locker.RLock()
	if item, ok := counter.FlowLimiterMap[hashName]; ok {
		item.LastTime = time.Now()
		counter.Locker.RUnlock()
		return item.Limter, nil
	}
	counter.Locker.RUnlock()

	item := &FlowLimiterItem{}
	var newLimiter Limiter
	if !isDistributed && ltype == 0 {
		newLimiter = rate.NewLimiter(rate.Limit(val), int(val*2))
	} else {
		capacity := val
		if ltype == DataTypeSecond {
			capacity = val * 2
		}
		newLimiter = NewDistributedLimiter(serverName, ltype, int64(val), int64(capacity))
	}

	item = &FlowLimiterItem{
		Name:          serverName,
		LType:         ltype,
		IsDistributed: isDistributed,
		Limter:        newLimiter,
		LastTime:      time.Now(),
	}

	counter.Locker.Lock()
	counter.FlowLimiterMap[hashName] = item
	defer counter.Locker.Unlock()
	return newLimiter, nil
}
