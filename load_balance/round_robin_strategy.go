package load_balance

import (
	"errors"
)

type RoundRobinStrategy struct {
	curIndex int
	rss      []string
}

func (r *RoundRobinStrategy) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RoundRobinStrategy) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss) //5
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RoundRobinStrategy) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinStrategy) GetAll() ([]string, error) {
	return r.rss, nil
}

func (r *RoundRobinStrategy) RemoveAll() error {
	r.rss = []string{}
	return nil
}

func init() {
	RegisterLoadBalanceStrategyHandler("round_robin", func() LoadBalanceStrategy {
		return &RoundRobinStrategy{}
	})
}
