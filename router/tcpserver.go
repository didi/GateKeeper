package router

import (
	"context"
	"encoding/json"
	"github.com/didichuxing/gatekeeper/dao"
	"github.com/didichuxing/gatekeeper/middleware"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/didichuxing/gatekeeper/tcpproxy"
	"log"
	"reflect"
	"sync"
	"time"
)

var (
	//TCPSrvHandler tcp服务回调
	TCPSrvHandler *TCPServer
)

//TCPServerRun 所有tcp服务启动
func TCPServerRun() {
	configer := service.SysConfMgr.GetModuleConfig()
	TCPSrvHandler = NewTCPServer(*configer)
	TCPSrvHandler.Run()
	go TCPSrvHandler.MonitorConfiger()
}

//TCPServerStop 所有tcp服务停止
func TCPServerStop() {
	TCPSrvHandler.Stop()
}

//TCPServer 全局的Tcp服务
type TCPServer struct {
	moduleServers map[string]*GateModuleTCPServer
	lock          sync.Mutex
	configer      dao.ModuleConfiger
}

//NewTCPServer tcpServer初始化
func NewTCPServer(configer dao.ModuleConfiger) *TCPServer {
	ss := map[string]*GateModuleTCPServer{}
	for _, tmpModule := range configer.Module {
		if tmpModule.Base.LoadType != "tcp" {
			continue
		}
		gs := NewGateModuleTCPServer(*tmpModule)
		ss[tmpModule.Base.Name] = gs
	}
	return &TCPServer{
		moduleServers: ss,
		lock:          sync.Mutex{},
		configer:      configer,
	}
}

//Run 运行多台Tcp Server
func (ts *TCPServer) Run() {
	for _, ms := range ts.moduleServers {
		go ms.Run()
	}
}

//Stop 停止多台Tcp Server
func (ts *TCPServer) Stop() {
	for _, ms := range ts.moduleServers {
		ms.Stop()
	}
}

//Reload 重新加载配置
func (ts *TCPServer) Reload(configer dao.ModuleConfiger) error {
	// 处理增加,和修改的
	for _, tmpModule := range configer.Module {
		if tmpModule.Base.LoadType != "tcp" {
			continue
		}
		if ms, ok := ts.moduleServers[tmpModule.Base.Name]; !ok {
			gs := NewGateModuleTCPServer(*tmpModule)
			go gs.Run() //不存在时则执行run
			ts.lock.Lock()
			ts.moduleServers[tmpModule.Base.Name] = gs
			ts.lock.Unlock()
		} else {
			oldModule, _ := json.Marshal(ms.Module)
			newModule, _ := json.Marshal(tmpModule)
			if string(oldModule) != string(newModule) {
				ms.Reload(*tmpModule) //存在且配置变化时执行reload
			}
			/*if !reflect.DeepEqual(ms.Module, tmpModule) {
				//oldModule,_:=json.Marshal(ms.Module)
				//newModule,_:=json.Marshal(tmpModule)
				//fmt.Println("old.Module",string(oldModule))
				//fmt.Println("new.Module",string(newModule))
				ms.Reload(*tmpModule) //存在且配置变化时执行reload
			}*/
		}
	}

	// 处理减少的
	for name, gs := range ts.moduleServers {
		module := configer.GetGateWayModule(name)
		if module == nil {
			gs.Stop()
			ts.lock.Lock()
			delete(ts.moduleServers, name)
			ts.lock.Unlock()
		}
	}
	return nil
}

//MonitorConfiger 配置检测
func (ts *TCPServer) MonitorConfiger() {
	for {
		select {
		case <-service.SysConfMgr.ConfigChangeNotice():
			configer := service.SysConfMgr.GetModuleConfig()
			if !reflect.DeepEqual(*configer, ts.configer) {
				ts.Reload(*configer)
			}
		}
	}
}

//GateModuleTCPServer 单个tcp服务器结构体
type GateModuleTCPServer struct {
	lock        sync.Mutex
	Module      dao.GatewayModule
	proxy       *tcpproxy.Proxy
	close       int
	avaliableIP []string
}

//NewGateModuleTCPServer 构造server
func NewGateModuleTCPServer(module dao.GatewayModule) *GateModuleTCPServer {
	return &GateModuleTCPServer{
		lock:        sync.Mutex{},
		Module:      module,
		proxy:       nil,
		close:       0,
		avaliableIP: []string{},
	}
}

//Run 服务启动，会阻塞的函数调用
func (s *GateModuleTCPServer) Run() error {
	defer func() {
		if err := recover(); err != nil {
			public.SysLogger.Error("Tcp Run Error:%v", err)
		}
	}()
	if err := public.CheckConnPort(s.Module.Base.FrontendAddr); err != nil {
		public.SysLogger.Error("Tcp Run Error:%v", err)
		log.Printf(" [WARN] TCPServer %s err:%v\n", s.Module.Base.FrontendAddr, err)
		return err
	}
	for {
		tcpNewProxyHandler := tcpproxy.NewProxy()
		s.avaliableIP = service.SysConfMgr.GetAvaliableIPList(s.Module.Base.Name)
		newIPWeightMap := service.SysConfMgr.GetConfIPWeightMap(&s.Module, public.IPDefaultWeight)

		//获取模块配置
		module := service.SysConfMgr.GetModuleConfigByName(s.Module.Base.Name)
		if module == nil {
			log.Printf(" [WARN] TCPServer %s removed\n", s.Module.Base.FrontendAddr)
			return nil
		}

		moduleIPList, cerr := service.SysConfMgr.GetModuleConfIPList(s.Module.Base.Name)
		if cerr != nil {
			log.Printf(" [WARN] TCPServer %s\n", s.Module.Base.FrontendAddr, cerr.Error())
			return nil
		}

		//如果可用ip列表为空，则兼容使用模块配置ip列表
		if len(s.avaliableIP) == 0 {
			s.avaliableIP = moduleIPList
		}

		//有关闭标记退出
		if s.close == 1 {
			log.Printf(" [WARN] TCPServer %s closed\n", s.Module.Base.FrontendAddr)
			return nil
		}

		//每个目标ip单独设定一个router及权重
		for _, ipAddr := range s.avaliableIP {
			w, ok := newIPWeightMap[ipAddr]
			weight := public.IPDefaultWeight
			if ok {
				weight = int(w)
			}
			tcpNewProxyHandler.AddWeightRoute(s.Module.Base.FrontendAddr, weight, tcpproxy.NewTCPRouter(
				&tcpproxy.DialProxy{
					Addr:        ipAddr,
					DialTimeout: time.Millisecond * time.Duration(s.Module.LoadBalance.ProxyConnectTimeout),
				}))
		}
		//设置中间件
		tcpNewProxyHandler.Use(
			s.Module.Base.FrontendAddr,
			middleware.TCPLimit(&s.Module))

		log.Printf(" [INFO] TCPServer %s listening\n", s.Module.Base.FrontendAddr)
		ctx, cancel := context.WithCancel(context.Background())
		go s.monitorAvaliableIP(ctx)
		s.proxy = tcpNewProxyHandler
		err := s.proxy.Run()
		if err != nil {
			//fmt.Println(s.Module.Base.Name+" in run err")
			cancel()
			log.Printf(" [WARN] TCPServer %s err:%v\n", s.Module.Base.FrontendAddr, err)
		}
	}
}

// 循环监控可用ip
func (s *GateModuleTCPServer) monitorAvaliableIP(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			public.CheckLogger.Error("%s monitorAvaiableIp:%v", s.Module.Base.Name, err)
		}
	}()
	t1 := time.NewTimer(time.Millisecond * time.Duration(s.Module.LoadBalance.CheckInterval))
	for {
		select {
		case <-t1.C:
			newIPList := service.SysConfMgr.GetAvaliableIPList(s.Module.Base.Name)
			//log.Printf("%s newIPList:%v\n",s.Module.Base.Name,newIPList)
			//模块ip确认变更了且ip大于0，则重启服务器
			if !reflect.DeepEqual(s.avaliableIP, newIPList) && len(newIPList) > 0 {
				s.lock.Lock()
				s.avaliableIP = newIPList
				s.lock.Unlock()
				s.Reload(s.Module)
				//log.Printf("%s reload oldIpList:%v, newIPList:%v \n",s.Module.Base.Name, s.avaliableIP,newIPList)
			}
			//log.Printf("%s monitorAvaliableIP:%v\n",s.Module.Base.Name,s.avaliableIP)
			t1.Reset(time.Millisecond * time.Duration(s.Module.LoadBalance.CheckInterval))
		case <-ctx.Done():
			//log.Println("monitorAvaliableIP close")
			return
		}
	}
}

//Stop 停止模块代理
func (s *GateModuleTCPServer) Stop() {
	s.lock.Lock()
	s.close = 1
	s.lock.Unlock()
	if s.proxy != nil {
		//fmt.Println(s.Module.Base.Name+" in stop")
		s.proxy.Close()
	}
}

//Restart 停止模块代理
func (s *GateModuleTCPServer) Restart() {
	if s.proxy != nil {
		//fmt.Println(s.Module.Base.Name+" in restart")
		s.proxy.Close()
		//fmt.Println(s.Module.Base.Name+" after restart")
	}
}

//Reload 重新加载模块
func (s *GateModuleTCPServer) Reload(module dao.GatewayModule) error {
	s.lock.Lock()
	s.Module = module
	s.lock.Unlock()
	s.Restart()
	return nil
}
