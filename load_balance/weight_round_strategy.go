package load_balance

import (
	"errors"
	"strconv"
)

type WeightRoundRobinStrategy struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
	conf     LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
}

func (r *WeightRoundRobinStrategy) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinStrategy) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		total += w.effectiveWeight
		w.currentWeight += w.effectiveWeight
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	best.currentWeight -= total
	return best.addr
}

func (r *WeightRoundRobinStrategy) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *WeightRoundRobinStrategy) GetAll() ([]string, error) {
	iplist := []string{}
	for _, item := range r.rss {
		iplist = append(iplist, item.addr)
	}
	return iplist, nil
}

func (r *WeightRoundRobinStrategy) RemoveAll() error {
	r.rss = []*WeightNode{}
	r.rsw = []int{}
	return nil
}

func init() {
	RegisterLoadBalanceStrategyHandler("weight_round", func() LoadBalanceStrategy {
		return &WeightRoundRobinStrategy{}
	})
}
