package public

import (
	"context"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/garyburd/redigo/redis"
	"sync"
	"sync/atomic"
	"time"
)

//FlowCounterHandler 全局流量统计回调
var FlowCounterHandler *FlowCounter

//FlowCounter 全局流量统计结构体
type FlowCounter struct {
	requestCountMap     map[string]*RequestCountService
	requestCountMapLock sync.RWMutex
	appCountMap         map[string]*APPCountService
	appCountMapLock     sync.RWMutex
}

//NewFlowCounter 创建FlowCounter
func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		requestCountMap:     make(map[string]*RequestCountService),
		requestCountMapLock: sync.RWMutex{},
		appCountMap:         make(map[string]*APPCountService),
		appCountMapLock:     sync.RWMutex{},
	}
}

//GetRequestCounter 获取一个模块统计，不存在就创建一个
func (c *FlowCounter) GetRequestCounter(moduleName string) *RequestCountService {
	c.requestCountMapLock.RLock()
	if counter, ok := c.requestCountMap[moduleName]; ok {
		c.requestCountMapLock.RUnlock()
		return counter
	}
	c.requestCountMapLock.RUnlock()
	c.requestCountMapLock.Lock()
	defer c.requestCountMapLock.Unlock()
	newCounter, err := NewRequestCountService(moduleName, 1*time.Second, 1)
	if err != nil {
		StatLogger.Error("GetRequestCounter_err:%v", err)
		return nil
	}
	c.requestCountMap[moduleName] = newCounter
	return newCounter
}

//GetAPPCounter 获取一个App统计不存在就创建一个
func (c *FlowCounter) GetAPPCounter(appID string) *APPCountService {
	c.appCountMapLock.RLock()
	if counter, ok := c.appCountMap[appID]; ok {
		c.appCountMapLock.RUnlock()
		return counter
	}
	c.appCountMapLock.RUnlock()
	c.appCountMapLock.Lock()
	defer c.appCountMapLock.Unlock()
	newCounter, err := NewAPPCountService(appID, 1*time.Second, 1)
	if err != nil {
		StatLogger.Error("GetAPPCounter_error:%v", err)
		return nil
	}
	c.appCountMap[appID] = newCounter
	return newCounter
}

//RequestCountService 请求计数结构体
type RequestCountService struct {
	ModuleName  string
	Interval    time.Duration
	Lock        sync.RWMutex
	ReqCount    int64
	TotalCount  int64
	QPS         int64
	Unix        int64
	TickerCount int64
	ReqDate     string
}

//NewRequestCountService 创建请求计数对象
func NewRequestCountService(moduleName string, interval time.Duration, maxCnt int) (*RequestCountService, error) {
	reqCounter := &RequestCountService{
		ModuleName:  moduleName,
		Interval:    interval,
		ReqCount:    0,
		QPS:         0,
		Unix:        0,
		TickerCount: 0,
		ReqDate:     "",
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				SysLogger.Error("NewRequestCountService_recover:%v", err)
			}
		}()
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据
			today := time.Now().In(TimeLocation).Format(DayFormat)
			redisKey := RequestModuleCounterPrefix + today + "_" + reqCounter.ModuleName
			todayhour := time.Now().In(TimeLocation).Format("2006010215")
			redisHourKey := RequestModuleHourCounterPrefix + todayhour + "_" + reqCounter.ModuleName

			RedisConfPipline(StatLogger, "default",
				func(c redis.Conn) {
					c.Send("INCRBY", redisKey, tickerCount)
					c.Send("EXPIRE", redisKey, 86400)
					c.Send("INCRBY", redisHourKey, tickerCount)
					c.Send("EXPIRE", redisHourKey, 86400)
				})
			if currentCount, err := redis.Int64(RedisConfDo(StatLogger, "default", "GET", redisKey)); err == nil {
				nowUnix := time.Now().Unix()
				nowDate := time.Now().In(lib.TimeLocation).Format(lib.DateFormat)
				if reqCounter.ReqDate != nowDate {
					reqCounter.ReqDate = nowDate
					reqCounter.TotalCount = 1
				}
				if reqCounter.Unix == 0 {
					reqCounter.Unix = time.Now().Unix()
				} else {
					if currentCount >= reqCounter.TotalCount && nowUnix > reqCounter.Unix {
						reqCounter.QPS = (currentCount - reqCounter.TotalCount) / (nowUnix - reqCounter.Unix)
						reqCounter.TotalCount = currentCount
						reqCounter.Unix = time.Now().Unix()
					}
				}
			}
		}
	}()
	return reqCounter, nil
}

//Increase 增加一次请求
func (o *RequestCountService) Increase(ctx context.Context, remoteAddr string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				SysLogger.Error("Increase_recover:%v", err)
			}
		}()
		atomic.AddInt64(&o.TickerCount, 1)
	}()
}

//GetHourCount 获取小时统计
func (o *RequestCountService) GetHourCount(dayhour string) (int64, error) {
	redisKey := RequestModuleHourCounterPrefix + dayhour + "_" + o.ModuleName
	return redis.Int64(lib.RedisConfDo(GetTraceContext(context.Background()), "default", "GET", redisKey))
}

//GetDayCount 获取天级统计
func (o *RequestCountService) GetDayCount(day string) (int64, error) {
	redisKey := RequestModuleCounterPrefix + day + "_" + o.ModuleName
	return redis.Int64(lib.RedisConfDo(GetTraceContext(context.Background()), "default", "GET", redisKey))
}

//APPCountService app统计结构体
type APPCountService struct {
	AppID       string
	Interval    time.Duration
	Lock        sync.RWMutex
	ReqCount    int64
	TotalCount  int64
	QPS         int64
	Unix        int64
	TickerCount int64
	ReqDate     string
}

//NewAPPCountService 创建统计结构体
func NewAPPCountService(appID string, interval time.Duration, maxCnt int) (*APPCountService, error) {
	reqCounter := &APPCountService{
		AppID:       appID,
		Interval:    interval,
		ReqCount:    0,
		QPS:         0,
		Unix:        0,
		TickerCount: 0,
		ReqDate:     "",
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				SysLogger.Error("NewAPPCountService_recover:%v", err)
			}
		}()
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据

			today := time.Now().In(TimeLocation).Format(DayFormat)
			totalAppKey := fmt.Sprintf("%s%s_%s", AccessControlAppIDTotalCallPrefix, today, appID)

			todayhour := time.Now().In(TimeLocation).Format("2006010215")
			redisHourKey := fmt.Sprintf("%s%s_%s", AccessControlAppIDHourTotalCallPrefix, todayhour, appID)
			RedisConfPipline(StatLogger, "default",
				func(c redis.Conn) {
					c.Send("INCRBY", totalAppKey, tickerCount)
					c.Send("EXPIRE", totalAppKey, 86400)
					c.Send("INCRBY", redisHourKey, tickerCount)
					c.Send("EXPIRE", redisHourKey, 86400)
				})

			if currentCount, err := redis.Int64(RedisConfDo(StatLogger, "default", "GET", totalAppKey)); err == nil {
				nowUnix := time.Now().Unix()
				nowDate := time.Now().In(lib.TimeLocation).Format(lib.DateFormat)
				if reqCounter.ReqDate != nowDate {
					reqCounter.ReqDate = nowDate
					reqCounter.TotalCount = 1
				}
				if reqCounter.Unix == 0 {
					reqCounter.Unix = time.Now().Unix()
				} else {
					if currentCount >= reqCounter.TotalCount && nowUnix > reqCounter.Unix {
						reqCounter.QPS = (currentCount - reqCounter.TotalCount) / (nowUnix - reqCounter.Unix)
						reqCounter.TotalCount = currentCount
						reqCounter.Unix = time.Now().Unix()
					}
				}
			}
		}
	}()
	return reqCounter, nil
}

//GetHourCount 获取小时统计
func (o *APPCountService) GetHourCount(dayhour string) (int64, error) {
	redisKey := AccessControlAppIDHourTotalCallPrefix + dayhour + "_" + o.AppID
	return redis.Int64(lib.RedisConfDo(GetTraceContext(context.Background()), "default", "GET", redisKey))
}

//GetDayCount 获取天级统计
func (o *APPCountService) GetDayCount(day string) (int64, error) {
	redisKey := AccessControlAppIDTotalCallPrefix + day + "_" + o.AppID
	return redis.Int64(lib.RedisConfDo(GetTraceContext(context.Background()), "default", "GET", redisKey))
}

//Increase 增加一次统计
func (o *APPCountService) Increase(context context.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				SysLogger.Error("Increase_recover:%v", err)
			}
		}()
		atomic.AddInt64(&o.TickerCount, 1)
	}()
}
