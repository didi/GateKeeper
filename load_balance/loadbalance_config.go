package load_balance

import "github.com/didi/gatekeeper/model"

type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
	CloseWatch()
}

type Observer interface {
	Update()
}

type CheckConfigHandler func(service *model.ServiceDetail) (LoadBalanceConf, error)

var CheckConfigHandlerMap map[string]CheckConfigHandler

func RegisterCheckConfigHandler(lbtype string, conf CheckConfigHandler) {
	if CheckConfigHandlerMap == nil {
		CheckConfigHandlerMap = map[string]CheckConfigHandler{}
	}
	CheckConfigHandlerMap[lbtype] = conf
}

func GetCheckConfigHandler(lbtype string) CheckConfigHandler {
	if CheckConfigHandlerMap == nil {
		return nil
	}
	handler, ok := CheckConfigHandlerMap[lbtype]
	if !ok {
		return CheckConfigHandlerMap["upstream_config"]
	}
	return handler
}
