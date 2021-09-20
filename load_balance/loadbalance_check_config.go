package load_balance

import (
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/model"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

type LoadBalanceCheckConf struct {
	observers    []Observer
	confIpWeight map[string]string
	activeList   []string
	format       string //单条数据格式 http://%s，%s方便替换成ip地址
	name         string
	closeChan    chan bool
}

func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

func (s *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

func (s *LoadBalanceCheckConf) CloseWatch() {
	s.closeChan <- true
	close(s.closeChan)
}

func (s *LoadBalanceCheckConf) WatchConf() {
	go func() {
		confIpErrNum := map[string]int{}
		log.Info().Msg(lib.Purple(fmt.Sprintf("checking [%s] config_list:%v active_list:%v\n", s.name, s.confIpWeight, s.activeList)))
	OUTFOR:
		for {
			select {
			case <-s.closeChan:
				break OUTFOR
			default:
			}
			changedList := []string{}
			for rs, _ := range s.confIpWeight {
				conn, err := net.DialTimeout("tcp", rs, time.Duration(DefaultCheckTimeout)*time.Second)
				if err == nil {
					conn.Close()
					if _, ok := confIpErrNum[rs]; ok {
						confIpErrNum[rs] = 0
					}
				}
				if err != nil {
					if _, ok := confIpErrNum[rs]; ok {
						confIpErrNum[rs] += 1
					} else {
						confIpErrNum[rs] = 1
					}
				}
				if confIpErrNum[rs] < DefaultCheckMaxErrNum {
					changedList = append(changedList, rs)
				}
			}
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			if !reflect.DeepEqual(changedList, s.activeList) {
				log.Info().Msg(lib.Purple(fmt.Sprintf("checking [%s] changed config_list:%v changed_list:%v\n", s.name, s.confIpWeight, changedList)))
				s.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceCheckConf(service *model.ServiceDetail) (LoadBalanceConf, error) {
	tmpstring := service.PluginConf.GetPath("upstream_config", "upstream_list").MustString()
	upConf, err := model.GetUpstreamConfigFromString(tmpstring)
	if err != nil {
		return nil, err
	}
	mConf := &LoadBalanceCheckConf{
		name:         service.Info.ServiceName,
		format:       fmt.Sprintf("%s%s", upConf.Schema, "%s"),
		activeList:   upConf.IpList,
		confIpWeight: upConf.IpWeight,
		closeChan:    make(chan bool, 1)}
	mConf.WatchConf()
	return mConf, nil
}

func init() {
	RegisterCheckConfigHandler("upstream_config", NewLoadBalanceCheckConf)
}
