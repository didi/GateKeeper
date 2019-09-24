package public

//RR 基于 权重round robin算法的接口
type RR interface {
	Next() interface{}
	Add(node interface{}, weight int)
	RemoveAll()
	Reset()
}

const (
	//RRNginx Nginx算法
	RRNginx = 0
	//RRLvs LVS算法
	RRLvs = 1
)

//NewWeightedRR 算法实现工厂类
func NewWeightedRR(rtype int) RR {
	if rtype == RRNginx {
		return &WNGINX{}
	} else if rtype == RRLvs {
		return &WLVS{}
	}
	return nil
}

//WeightNginx 节点结构
type WeightNginx struct {
	Node            interface{}
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func (ww *WeightNginx) fail() {
	ww.EffectiveWeight -= ww.Weight
	if ww.EffectiveWeight < 0 {
		ww.EffectiveWeight = 0
	}
}

//WNGINX nginx算法实现类
type WNGINX struct {
	nodes []*WeightNginx
	n     int
}

//Add 增加权重节点
func (w *WNGINX) Add(node interface{}, weight int) {
	weighted := &WeightNginx{
		Node:            node,
		Weight:          weight,
		EffectiveWeight: weight}
	w.nodes = append(w.nodes, weighted)
	w.n++
}

//RemoveAll 移除所有节点
func (w *WNGINX) RemoveAll() {
	w.nodes = w.nodes[:0]
	w.n = 0
}

//Next 下次轮询事件
func (w *WNGINX) Next() interface{} {
	if w.n == 0 {
		return nil
	}
	if w.n == 1 {
		return w.nodes[0].Node
	}

	return nextWeightedNode(w.nodes).Node
}

func nextWeightedNode(nodes []*WeightNginx) (best *WeightNginx) {
	total := 0

	for i := 0; i < len(nodes); i++ {
		w := nodes[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	if best == nil {
		return nil
	}
	best.CurrentWeight -= total
	return best
}

//Reset 重置权重
func (w *WNGINX) Reset() {
	for _, s := range w.nodes {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0
	}
}

//WeightLvs 节点结构
type WeightLvs struct {
	Node   interface{}
	Weight int
}

//WLVS lvs算法实现类
type WLVS struct {
	nodes []*WeightLvs
	n     int
	gcd   int //通用的权重因子
	maxW  int //最大权重
	i     int //被选择的次数
	cw    int //当前的权重值
}

//Next 下次轮询事件
func (w *WLVS) Next() interface{} {
	if w.n == 0 {
		return nil
	}

	if w.n == 1 {
		return w.nodes[0].Node
	}

	for {
		w.i = (w.i + 1) % w.n
		if w.i == 0 {
			w.cw = w.cw - w.gcd
			if w.cw <= 0 {
				w.cw = w.maxW
				if w.cw == 0 {
					return nil
				}
			}
		}
		if w.nodes[w.i].Weight >= w.cw {
			return w.nodes[w.i].Node
		}
	}
}

//Add 增加权重节点
func (w *WLVS) Add(node interface{}, weight int) {
	weighted := &WeightLvs{Node: node, Weight: weight}
	if weight > 0 {
		if w.gcd == 0 {
			w.gcd = weight
			w.maxW = weight
			w.i = -1
			w.cw = 0
		} else {
			w.gcd = gcd(w.gcd, weight)
			if w.maxW < weight {
				w.maxW = weight
			}
		}
	}
	w.nodes = append(w.nodes, weighted)
	w.n++
}

func gcd(x, y int) int {
	var t int
	for {
		t = (x % y)
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}

//RemoveAll 移除所有节点
func (w *WLVS) RemoveAll() {
	w.nodes = w.nodes[:0]
	w.n = 0
	w.gcd = 0
	w.maxW = 0
	w.i = -1
	w.cw = 0
}

//Reset 重置权重
func (w *WLVS) Reset() {
	w.i = -1
	w.cw = 0
}
