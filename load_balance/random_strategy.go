package load_balance

import (
	"errors"
	"math/rand"
)

type RandomStrategy struct {
	curIndex int
	rss      []string
}

func (r *RandomStrategy) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RandomStrategy) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

func (r *RandomStrategy) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomStrategy) GetAll() ([]string, error) {
	return r.rss, nil
}

func (r *RandomStrategy) RemoveAll() error {
	r.rss = []string{}
	return nil
}

func init() {
	RegisterLoadBalanceStrategyHandler("random", func() LoadBalanceStrategy {
		return &RandomStrategy{}
	})
}
