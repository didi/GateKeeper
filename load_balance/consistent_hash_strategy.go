package load_balance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashStrategy struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int               //复制因子
	keys     UInt32Slice       //已排序的节点hash切片
	hashMap  map[uint32]string //节点哈希和Key的map,键是hash值，值是节点key
}

func NewConsistentHashStrategy(replicas int, fn Hash) *ConsistentHashStrategy {
	m := &ConsistentHashStrategy{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (c *ConsistentHashStrategy) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashStrategy) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	c.mux.Lock()
	defer c.mux.Unlock()
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}
	sort.Sort(c.keys)
	return nil
}

func (c *ConsistentHashStrategy) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}
	hash := c.hash([]byte(key))
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })
	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}

func (r *ConsistentHashStrategy) GetAll() ([]string, error) {
	iplist := []string{}
	for _, item := range r.hashMap {
		iplist = append(iplist, item)
	}
	return iplist, nil
}

func (c *ConsistentHashStrategy) RemoveAll() error {
	c.keys = UInt32Slice{}
	c.hashMap = map[uint32]string{}
	return nil
}

func init() {
	RegisterLoadBalanceStrategyHandler("consistent_hash", func() LoadBalanceStrategy {
		return NewConsistentHashStrategy(32, nil)
	})
}
