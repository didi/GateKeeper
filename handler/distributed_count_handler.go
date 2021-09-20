package handler

import (
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/public"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync/atomic"
	"time"
)

type DistributedCountService struct {
	Name        string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
	closeChan   chan bool
}

func NewDistributedCountService(appID string, interval time.Duration) *DistributedCountService {
	reqCounter := &DistributedCountService{
		Name:     appID,
		Interval: interval,
		QPS:      0,
		Unix:     0,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		ticker := time.NewTicker(interval)
	OUTFOR:
		for {
			select {
			case <-reqCounter.closeChan:
				continue OUTFOR
			case <-ticker.C:
				tickerCount := atomic.LoadInt64(&reqCounter.TickerCount)
				atomic.StoreInt64(&reqCounter.TickerCount, 0)
				currentTime := time.Now()
				dayKey := reqCounter.GetDayKey(currentTime)
				hourKey := reqCounter.GetHourKey(currentTime)

				//todo 修改成odin上报或promius上报？
				//次时tickerCount是单机流量统计
				if tickerCount > 0 {
					if err := public.RedisConfPipline(func(c redis.Conn) {
						c.Send("INCRBY", dayKey, tickerCount)
						c.Send("EXPIRE", dayKey, 86400*2)
						c.Send("INCRBY", hourKey, tickerCount)
						c.Send("EXPIRE", hourKey, 86400*2)
					}); err != nil {
						fmt.Println("RedisConfPipline err", err)
						continue
					}
				}
				totalCount, err := reqCounter.GetDayData(currentTime)
				if err != nil && err.Error() != "redigo: nil returned" {
					fmt.Println("reqCounter.GetDayData err", err)
					continue
				}
				nowUnix := time.Now().Unix()
				if reqCounter.Unix == 0 {
					reqCounter.Unix = time.Now().Unix()
					continue
				}
				tickerCount = totalCount - reqCounter.TotalCount
				//次时tickerCount是分布式流量统计
				if nowUnix > reqCounter.Unix {
					reqCounter.TotalCount = totalCount
					reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
					reqCounter.Unix = time.Now().Unix()
				}
			}
		}
	}()
	return reqCounter
}

func (o *DistributedCountService) Close() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
	}()
	if o.closeChan != nil {
		close(o.closeChan)
	}
}

func (o *DistributedCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(lib.TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", public.RedisFlowDayKey, dayStr, o.Name)
}

func (o *DistributedCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", public.RedisFlowHourKey, hourStr, o.Name)
}

func (o *DistributedCountService) GetHourData(t time.Time) (int64, error) {
	return redis.Int64(public.RedisConfDo("GET", o.GetHourKey(t)))
}

func (o *DistributedCountService) GetDayData(t time.Time) (int64, error) {
	return redis.Int64(public.RedisConfDo("GET", o.GetDayKey(t)))
}

func (o *DistributedCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		atomic.AddInt64(&o.TickerCount, 1)
	}()
}
