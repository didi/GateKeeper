package public

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

//FlowLimiterHandler 全局流量统计回调
var FlowLimiterHandler *FlowLimiter

//LimitVisitor 流控结构体
type LimitVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

//FlowLimiter 流控管理器
type FlowLimiter struct{
	limitVisitors map[string]*LimitVisitor
	limitLock     sync.RWMutex
}

//NewFlowLimiter 创建对象
func NewFlowLimiter() *FlowLimiter {
	t:=&FlowLimiter{
		limitVisitors:make(map[string]*LimitVisitor),
		limitLock:sync.RWMutex{},
	}
	go t.CleanupLimitVisitors()
	return t
}

//AddAppVisitor 创建app流控
func(t *FlowLimiter) AddAppVisitor(appID string,qps int64) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	t.limitLock.Lock()
	t.limitVisitors[appID] = &LimitVisitor{limiter, time.Now()}
	t.limitLock.Unlock()
	return limiter
}

//AddLimitVisitor 创建流控对象
func(t *FlowLimiter) AddLimitVisitor(name string, qps int64) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	t.limitLock.Lock()
	t.limitVisitors[name] = &LimitVisitor{limiter, time.Now()}
	t.limitLock.Unlock()
	return limiter
}

//GetAPPLimitVisitor 获取app流控对象,不存在就创建
func(t *FlowLimiter) GetAPPLimitVisitor(appID string,qps int64) *rate.Limiter {
	t.limitLock.RLock()
	v, exists := t.limitVisitors[appID]
	if !exists {
		t.limitLock.RUnlock()
		return t.AddAppVisitor(appID,qps)
	}
	v.lastSeen = time.Now()
	t.limitLock.RUnlock()
	return v.limiter
}

//GetModuleIPVisitor 获取module流控对象,不存在就创建
func(t *FlowLimiter) GetModuleIPVisitor(moduleIPAddr string, qps int64) *rate.Limiter {
	t.limitLock.RLock()
	v, exists := t.limitVisitors[moduleIPAddr]
	if !exists {
		t.limitLock.RUnlock()
		return t.AddLimitVisitor(moduleIPAddr, qps)
	}
	v.lastSeen = time.Now()
	t.limitLock.RUnlock()
	return v.limiter
}

//CleanupLimitVisitors 定时清空流控对象
func(t *FlowLimiter) CleanupLimitVisitors() {
	for {
		time.Sleep(time.Minute)
		t.limitLock.Lock()
		for ip, v := range t.limitVisitors {
			if time.Now().Sub(v.lastSeen) > 5*time.Second {
				delete(t.limitVisitors, ip)
			}
		}
		t.limitLock.Unlock()
	}
}