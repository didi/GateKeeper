package handler

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/public"
	"github.com/garyburd/redigo/redis"
	"math"
	"time"
)

const (
	DataTypeSecond = iota
	DataTypeMinute
	DataTypeHour
)

type DistributedLimiter struct {
	Name     string
	Dtype    int //0=qps 1=qpm 2=qph
	Rate     int64
	Capacity int64
}

func NewDistributedLimiter(name string, dtype int, rate, capacity int64) *DistributedLimiter {
	if dtype == DataTypeMinute {
		rate = capacity / 60
	}
	if dtype == DataTypeHour {
		rate = capacity / 3600
	}
	if rate < 1 {
		rate = 1
	}
	return &DistributedLimiter{
		Name:     name,
		Dtype:    dtype,
		Rate:     rate,
		Capacity: capacity,
	}
}

func RedisScript(script string, args ...interface{}) (interface{}, error) {
	c, err := lib.RedisConnFactory("default")
	if err != nil {
		return nil, err
	}
	defer c.Close()
	lua := redis.NewScript(1, script)
	return lua.Do(c, args...)
}

func (d *DistributedLimiter) Allow() bool {
	luaScript := `--$lua = <<<SCRIPT
local key = KEYS[1]               --每秒一个，如：sv1_1625898937
local limit = tonumber(ARGV[1])   --限流大小，如：20
local current = tonumber(redis.call('get', key) or "0")
if ( current + 1 > limit )
then
    return 0
else
    redis.call("INCRBY", key,"1") --自增
    redis.call("expire", key,"2") --2秒超时
    return 1
end
--SCRIPT;
`
	allow, err := RedisScript(string(luaScript), d.Name, d.Capacity)
	if err != nil {
		log.Printf("DistributedLimiter RedisScript Error: %v\n", err)
		return false
	}
	if allow.(int64) == 0 {
		return false
	}
	return true
}

func (d *DistributedLimiter) AllowV1() bool {
	fillTime := float64(d.Capacity) / float64(d.Rate)
	ttl := math.Floor(fillTime * 2)
	redisKey := public.DistributedLimiterPrefix + d.Name
	redisMap, _ := redis.Int64Map(public.RedisConfDo("HGETALL", redisKey))
	lastTokens, ok := redisMap["tokens"]
	if !ok {
		lastTokens = lastTokens
	}
	lastRefreshed, ok := redisMap["timestamp"]
	if !ok {
		lastRefreshed = 0
	}
	delta := math.Max(0, float64(time.Now().Unix()-lastRefreshed))
	filledTokens := math.Min(float64(d.Capacity), float64(lastTokens)+(delta*float64(d.Rate)))
	allowed := false
	newTokens := filledTokens
	if filledTokens >= 1 {
		allowed = true
		newTokens = filledTokens - 1
	}
	public.RedisConfDo("HMSET", redisKey, "tokens", newTokens, "timestamp")
	public.RedisConfDo("EXPIRE", redisKey, ttl)
	return allowed
}
