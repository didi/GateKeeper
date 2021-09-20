package load_balance

import (
	"strings"
)

type LoadBalanceStrategy interface {
	Add(...string) error
	RemoveAll() error
	GetAll() ([]string, error)
	Get(string) (string, error)
}

type LoadBalance struct {
	strategy LoadBalanceStrategy
	conf     LoadBalanceConf
}

func (r *LoadBalance) Add(params ...string) error {
	r.strategy.Add(params...)
	return nil
}

func (r *LoadBalance) Get(params string) (string, error) {
	return r.strategy.Get(params)
}

func (r *LoadBalance) GetAll() ([]string, error) {
	return r.strategy.GetAll()
}

func (r *LoadBalance) Update() {
	r.strategy.RemoveAll()
	for _, ip := range r.conf.GetConf() {
		r.strategy.Add(strings.Split(ip, ",")...)
	}
}

func (r *LoadBalance) Close() {
	r.conf.CloseWatch()
}

func NewLoadBalance(strategy LoadBalanceStrategy, conf LoadBalanceConf) *LoadBalance {
	return &LoadBalance{
		strategy: strategy,
		conf:     conf,
	}
}

func LoadBanlanceFactorWithStrategy(strategy LoadBalanceStrategy, conf LoadBalanceConf) *LoadBalance {
	lb := NewLoadBalance(strategy, conf)
	conf.Attach(lb)
	lb.Update()
	return lb
}

type LoadBalanceStrategyHandler func() LoadBalanceStrategy

var LoadBalanceStrategyHandlerMap map[string]LoadBalanceStrategyHandler

func RegisterLoadBalanceStrategyHandler(name string, handler LoadBalanceStrategyHandler) {
	if LoadBalanceStrategyHandlerMap == nil {
		LoadBalanceStrategyHandlerMap = map[string]LoadBalanceStrategyHandler{}
	}
	LoadBalanceStrategyHandlerMap[name] = handler
}

func GetLoadBalanceStrategy(name string) LoadBalanceStrategy {
	if LoadBalanceStrategyHandlerMap == nil {
		return nil
	}
	handler := LoadBalanceStrategyHandlerMap[name]
	return handler()
}
