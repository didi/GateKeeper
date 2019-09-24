package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/didichuxing/gatekeeper/dao"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/BurntSushi/toml"
	"net/http"
	"reflect"
	"strconv"

	"github.com/e421083458/golang_common/lib"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"
)

//SysConfigManage 系统配置管理器
type SysConfigManage struct {
	moduleConfig       *dao.ModuleConfiger
	moduleConfigLocker sync.RWMutex

	moduleIPListMap       map[string][]string //可用ip
	moduleIPListMapLocker sync.RWMutex

	moduleActiveIPListMap       map[string][]string //活动ip
	moduleActiveIPListMapLocker sync.RWMutex

	moduleForbidIPListMap       map[string][]string //禁用ip
	moduleForbidIPListMapLocker sync.RWMutex

	moduleProxyFuncMap       map[string]func(rr public.RR) *httputil.ReverseProxy
	moduleProxyFuncMapLocker sync.RWMutex

	moduleRRMap       map[string]public.RR
	moduleRRMapLocker sync.RWMutex

	appConfig       *dao.APPConfiger
	appConfigLocker sync.RWMutex

	loadConfigContext context.Context //重新载入配置时，需要执行close
	loadConfigCancel  func()          //停止配置自动检查
}

//SysConfMgr 全局系统配置变量
var SysConfMgr *SysConfigManage

//NewSysConfigManage 实例化全局系统配置
func NewSysConfigManage() *SysConfigManage {
	return &SysConfigManage{
		moduleIPListMap:       map[string][]string{},
		moduleActiveIPListMap: map[string][]string{},
		moduleForbidIPListMap: map[string][]string{},
		moduleProxyFuncMap:    map[string]func(rr public.RR) *httputil.ReverseProxy{},
		moduleRRMap:           map[string]public.RR{},
	}
}

//InitConfig 初始化配置
func (s *SysConfigManage) InitConfig() {
	s.loadConfigContext, s.loadConfigCancel = context.WithCancel(context.Background())
	if err := s.refreshAPPConfig(); err != nil {
		panic(err)
	}
	if err := s.refreshModuleConfig(); err != nil {
		panic(err)
	}
	s.checkIPList()
	s.configModuleRR()
	s.configModuleProxyMap()
}

//ReloadConfig 刷新配置
func (s *SysConfigManage) ReloadConfig() {
	//刷新获取配置
	s.refreshAPPConfig()
	s.refreshModuleConfig()

	//如果非首次加载，执行cancel
	if s.loadConfigCancel != nil {
		s.loadConfigCancel()
	}
	s.loadConfigContext, s.loadConfigCancel = context.WithCancel(context.Background())

	//检测及配置
	s.checkIPList()
	s.configModuleRR()
	s.configModuleProxyMap()
}

//MonitorConfig 自动刷新配置
func (s *SysConfigManage) MonitorConfig() {
	//定期刷新模块配置，确保新模块加载成功
	loadInterval := lib.GetIntConf("base.base.conf_load_interval")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				public.SysLogger.Error("InternalRefresh_recover:%v", err)
			}
		}()
		for {
			time.Sleep(time.Duration(loadInterval) * time.Millisecond)
			s.ReloadConfig()
		}
	}()
}

//ConfigChangeNotice 配置变更通知
func (s *SysConfigManage) ConfigChangeNotice() (<-chan struct{}) {
	return s.loadConfigContext.Done()
}

//GetModuleConfig 获取全局系统配置
func (s *SysConfigManage) GetModuleConfig() *dao.ModuleConfiger {
	s.moduleConfigLocker.RLock()
	defer s.moduleConfigLocker.RUnlock()
	return s.moduleConfig
}

//GetModuleConfigByName 通过模块名获取模块配置
func (s *SysConfigManage) GetModuleConfigByName(name string) *dao.GatewayModule {
	for _, item := range s.GetModuleConfig().Module {
		if item.Base.Name == name {
			return item
		}
	}
	return nil
}

//GetActiveIPList 获取活动IP
func (s *SysConfigManage) GetActiveIPList(moduleName string) []string {
	s.moduleActiveIPListMapLocker.RLock()
	activeIPList, _ := s.moduleActiveIPListMap[moduleName]
	s.moduleActiveIPListMapLocker.RUnlock()
	return activeIPList
}

//GetForbidIPList 获取禁用IP
func (s *SysConfigManage) GetForbidIPList(moduleName string) []string {
	s.moduleForbidIPListMapLocker.RLock()
	ipList, _ := s.moduleForbidIPListMap[moduleName]
	s.moduleForbidIPListMapLocker.RUnlock()
	return ipList
}

//GetAvaliableIPList 获取可用IP
func (s *SysConfigManage) GetAvaliableIPList(moduleName string) []string {
	s.moduleIPListMapLocker.RLock()
	ipList, _ := s.moduleIPListMap[moduleName]
	s.moduleIPListMapLocker.RUnlock()
	return ipList
}

//GetModuleIPList 获取当前可用的ip列表
func (s *SysConfigManage) GetModuleIPList(moduleName string) ([]string, error) {
	ipList := s.GetAvaliableIPList(moduleName)
	if len(ipList) != 0 {
		return ipList, nil
	}
	return []string{}, errors.New("module ip empty")
}

//GetModuleConfIPList 获取配置的ip列表
func (s *SysConfigManage) GetModuleConfIPList(moduleName string) ([]string, error) {
	if moduleConf := s.GetModuleConfigByName(moduleName); moduleConf != nil {
		ipList := strings.Split(moduleConf.LoadBalance.IPList, ",")
		return ipList, nil
	}
	return []string{}, errors.New("module ip empty")
}

//GetModuleRR 获取模块的负载信息
func (s *SysConfigManage) GetModuleRR(moduleName string) (public.RR, error) {
	s.moduleRRMapLocker.RLock()
	defer s.moduleRRMapLocker.RUnlock()
	rr, ok := s.moduleRRMap[moduleName]
	if ok {
		return rr, nil
	}
	return nil, errors.New("module rr empty")
}

//GetConfIPWeightMap 返回模块对应的ip及权重
func (s *SysConfigManage) GetConfIPWeightMap(module *dao.GatewayModule, defaultWeight int64) map[string]int64 {
	confIPList := strings.Split(module.LoadBalance.IPList, ",")
	confWeightList := strings.Split(module.LoadBalance.WeightList, ",")
	confIPWeightMap := map[string]int64{}
	for index, ipAddr := range confIPList {
		if len(confWeightList) >= index+1 {
			weight, err := strconv.ParseInt(confWeightList[index], 10, 64)
			if err != nil {
				weight = defaultWeight
			}
			confIPWeightMap[ipAddr] = weight
		} else {
			confIPWeightMap[ipAddr] = defaultWeight
		}
	}
	return confIPWeightMap
}

//GetAppConfigByAPPID 获取租户数据
func (s *SysConfigManage) GetAppConfigByAPPID(appid string) (*dao.GatewayAPP, error) {
	for _, config := range s.appConfig.AuthirizedApps {
		if config.AppID == appid {
			return config, nil
		}
	}
	return nil, errors.New("app_id_is_empty")
}

//GetModuleHTTPProxy 获取http代理方法
func (s *SysConfigManage) GetModuleHTTPProxy(moduleName string) (*httputil.ReverseProxy, error) {
	rr, err := s.GetModuleRR(moduleName)
	if err != nil {
		return nil, err
	}
	s.moduleProxyFuncMapLocker.RLock()
	defer s.moduleProxyFuncMapLocker.RUnlock()
	proxyFunc, ok := s.moduleProxyFuncMap[moduleName]
	if ok {
		return proxyFunc(rr), nil
	}
	return nil, errors.New("module proxy empty")
}

//refreshModuleConfig 刷新模块配置到内存
func (s *SysConfigManage) refreshModuleConfig() error {
	defer func() {
		if err := recover(); err != nil {
			public.SysLogger.Error("RefreshModule_recover:%v", err)
		}
	}()
	configFile := lib.GetConfFilePath("module.toml")
	public.SysLogger.Info("module_file:%s", configFile)
	fileConf, err := s.getFileModuleConf(configFile)
	if err != nil {
		public.SysLogger.Error("GetFileModuleConf_error:%v", err)
		return err
	}

	dbConf, err := s.getDBModuleConf(true)
	if err != nil {
		public.SysLogger.Error("GetDBModuleConf_error:%v", err)
	}
	public.SysLogger.Info("GetDBModuleConf:%v", dbConf)

	//如果db挂了默认降级走file
	if dbConf != nil {
		s.moduleConfigLocker.Lock()
		s.moduleConfig = dbConf
		s.moduleConfigLocker.Unlock()
		public.SysLogger.Info("module_configured_by_db.")
		err := s.writeFileModuleConf(configFile, dbConf)
		if err != nil {
			public.SysLogger.Error("WriteFileModuleConf_error:%v", err)
		} else {
			public.SysLogger.Info("module_file_was_override.")
		}
	} else if fileConf != nil {
		s.moduleConfigLocker.Lock()
		s.moduleConfig = fileConf
		s.moduleConfigLocker.Unlock()
		public.SysLogger.Info("module_configured_by_file.")
	} else {
		public.SysLogger.Info("get_dbConf_and_fileConf_both_error")
		return err
	}
	public.SysLogger.Info("ModuleConf:%v", s.GetModuleConfig())
	return nil
}

func (s *SysConfigManage) checkModuleIPList(balance *dao.GatewayLoadBalance) []string {
	newIPList := []string{}
	ipList := strings.Split(balance.IPList, ",")
	for _, ip := range ipList {
		checkURL := fmt.Sprintf("http://%s%s", ip, balance.CheckURL)
		if balance.CheckMethod == "httpchk" {
			response, _, err := public.HTTPGET(public.CheckLogger, checkURL, nil, balance.CheckInterval, nil)
			if err != nil || response.StatusCode != 200 {
				public.CheckLogger.Warn("dltag=httpchk_failure|url=%v", checkURL)
			} else {
				newIPList = append(newIPList, ip)
			}
		}
		if balance.CheckMethod == "tcpchk" {
			conn, err := net.DialTimeout("tcp", ip, time.Millisecond*time.Duration(int64(balance.CheckInterval)))
			if err != nil {
				public.CheckLogger.Warn("dltag=_com_tcp_failure|ip=%v", ip)
			} else {
				conn.Close()
				public.CheckLogger.Warn("dltag=_com_tcp_success|ip=%v", ip)
				newIPList = append(newIPList, ip)
			}
		}
	}
	return newIPList
}

//配置模块服务发现检测
//已运行模块周期刷新，直到IpContext停止
func (s *SysConfigManage) checkIPList() {
	moduleConfiger := s.GetModuleConfig()
	s.moduleIPListMapLocker.Lock()
	s.moduleForbidIPListMapLocker.Lock()
	for _, module := range moduleConfiger.Module {
		s.moduleIPListMap[module.Base.Name] = strings.Split(module.LoadBalance.IPList, ",")
		s.moduleForbidIPListMap[module.Base.Name] = strings.Split(module.LoadBalance.ForbidList, ",")
	}
	s.moduleForbidIPListMapLocker.Unlock()
	s.moduleIPListMapLocker.Unlock()
	for _, modulePt := range moduleConfiger.Module {
		module := modulePt
		go func() {
			defer func() {
				if err := recover(); err != nil {
					public.CheckLogger.Warn("checkModuleIpList_recover:%v", err)
				}
			}()
			t1 := time.NewTimer(time.Second * 0)
		Loop:
			for {
				select {
				case <-t1.C:
					activeIPList := s.checkModuleIPList(module.LoadBalance)
					s.moduleActiveIPListMapLocker.Lock()
					s.moduleActiveIPListMap[module.Base.Name] = activeIPList
					s.moduleActiveIPListMapLocker.Unlock()

					s.moduleForbidIPListMapLocker.Lock()
					forbidIPList, ok := s.moduleForbidIPListMap[module.Base.Name]
					s.moduleForbidIPListMapLocker.Unlock()
					if !ok {
						forbidIPList = []string{}
					}

					//剔除禁用节点
					newIPList := []string{}
					for _, newIP := range activeIPList {

						if !public.InStringList(newIP, forbidIPList) {
							newIPList = append(newIPList, newIP)
						}
					}

					configIPList := strings.Split(module.LoadBalance.IPList, ",")
					s.moduleIPListMapLocker.Lock()
					s.moduleIPListMap[module.Base.Name] = newIPList
					s.moduleIPListMapLocker.Unlock()
					public.CheckLogger.Info("%s CheckModuleIpList newIPList=%+v configIPList=%+v", module.Base.Name, newIPList, configIPList)
					t1.Reset(time.Millisecond * time.Duration(module.LoadBalance.CheckInterval))

				case <-s.ConfigChangeNotice():
					public.CheckLogger.Info(module.Base.Name + "_CheckModuleIpList done")
					break Loop
				}
			}
		}()
	}
}

func (s *SysConfigManage) checkModuleConf(conf *dao.ModuleConfiger) error {
	if conf == nil || len(conf.Module) == 0 {
		return errors.New("conf is empty")
	}
	for _, confItem := range conf.Module {
		if confItem.Base == nil {
			return errors.New("module.base is empty")
		}
		if confItem.LoadBalance == nil {
			return errors.New("module.load_balance is empty")
		}

		//validator
		errs := public.ValidatorHandler.Struct(confItem)
		if errs != nil {
			return errs
		}
	}
	return nil
}

func (s *SysConfigManage) getFileModuleConf(confPath string) (*dao.ModuleConfiger, error) {
	moduleConf := &dao.ModuleConfiger{}
	file, err := os.Open(confPath)
	if err != nil {
		return moduleConf, err
	}
	defer file.Close()
	bts, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(bts), moduleConf); err != nil {
		return nil, err
	}
	if err := s.checkModuleConf(moduleConf); err != nil {
		return nil, err
	}
	return moduleConf, nil
}

func (s *SysConfigManage) writeFileModuleConf(confPath string, moduleConf *dao.ModuleConfiger) error {
	var buffer bytes.Buffer
	e := toml.NewEncoder(&buffer)
	if err := e.Encode(moduleConf); err != nil {
		return err
	}
	bakPath := strings.Replace(confPath, ".toml", "_bak.toml", -1)
	os.Remove(bakPath)
	if err := os.Rename(confPath, bakPath); err != nil {
		return err
	}
	if err := ioutil.WriteFile(confPath, buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func (s *SysConfigManage) getDBModuleConf(ischeck bool) (*dao.ModuleConfiger, error) {
	defer func() {
		if err := recover(); err != nil {
			public.SysLogger.Error("GetDBModuleConf_recover:%v", err)
		}
	}()

	moduleConf := &dao.ModuleConfiger{Module: []*dao.GatewayModule{}}
	bases, err := (&dao.GatewayModuleBase{}).GetAll()
	if err != nil {
		return nil, err
	}
	matchRuleArr, err := (&dao.GatewayMatchRule{}).GetAll()
	if err != nil {
		return nil, err
	}
	accessControlArr, err := (&dao.GatewayAccessControl{}).GetAll()
	if err != nil {
		return nil, err
	}
	loadBalanceArr, err := (&dao.GatewayLoadBalance{}).GetAll()
	if err != nil {
		return nil, err
	}
	for _, base := range bases {
		matchRules := []*dao.GatewayMatchRule{}
		for _, x := range matchRuleArr {
			if x.ModuleID == base.ID {
				matchRules = append(matchRules, x)
			}
		}
		accessControl := &dao.GatewayAccessControl{}
		for _, x := range accessControlArr {
			if x.ModuleID == base.ID {
				accessControl = x
			}
		}
		loadBalance := &dao.GatewayLoadBalance{}
		for _, x := range loadBalanceArr {
			if x.ModuleID == base.ID {
				loadBalance = x
			}
		}
		// why check base != nil at the end, and repeat with CheckModuleConf
		if base != nil && loadBalance != nil {
			moduleConf.Module = append(moduleConf.Module, &dao.GatewayModule{
				Base:          base,
				MatchRule:     matchRules,
				//DataFilter:    dataFilter,
				AccessControl: accessControl,
				LoadBalance:   loadBalance,
			})
		}
	}
	if ischeck {
		if err := s.checkModuleConf(moduleConf); err != nil {
			return nil, err
		}
	}
	return moduleConf, nil
}

//刷新租户信息到内存
func (s *SysConfigManage) refreshAPPConfig() error {
	defer func() {
		if err := recover(); err != nil {
			public.SysLogger.Error("RefreshAPP_recover:%v", err)
		}
	}()
	configFile := lib.GetConfFilePath("app.toml")
	public.SysLogger.Info("module_file:%s", configFile)
	fileConf, err := s.getFileAPPConf(configFile)
	if err != nil {
		public.SysLogger.Error("GetFileAPPConf_error:%v", err)
		return err
	}
	public.SysLogger.Info("GetFileAPPConf:%v", fileConf)
	dbConf, err := s.getDBAPPConf(true)
	if err != nil {
		public.SysLogger.Error("GetDBAPPConf_error:%v", err)
	}
	public.SysLogger.Info("GetDBAPPConf:%v", dbConf)

	//如果db挂了默认降级走file
	if dbConf != nil {
		s.appConfigLocker.Lock()
		s.appConfig = dbConf
		s.appConfigLocker.Unlock()
		public.SysLogger.Info("app_configured_by_db.")
		err := s.writeFileAppConf(configFile, dbConf)
		if err != nil {
			public.SysLogger.Error("WriteFileModuleConf_error:%v", err)
		} else {
			public.SysLogger.Info("app_file_was_override.")
		}
	} else if fileConf != nil {
		s.appConfigLocker.Lock()
		s.appConfig = fileConf
		s.appConfigLocker.Unlock()
		public.SysLogger.Info("app_configured_by_file.")
	} else {
		public.SysLogger.Info("get_dbConf_and_fileConf_both_error")
		return err
	}
	public.SysLogger.Info("APPConf:%v", s.appConfig)
	return nil
}

func (s *SysConfigManage) getFileAPPConf(confPath string) (*dao.APPConfiger, error) {
	appConf := &dao.APPConfiger{}
	file, err := os.Open(confPath)
	if err != nil {
		return appConf, err
	}
	defer file.Close()
	bts, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(bts), appConf); err != nil {
		return nil, err
	}
	if err := s.checkAppConf(appConf); err != nil {
		return nil, err
	}
	return appConf, nil
}

func (s *SysConfigManage) checkAppConf(conf *dao.APPConfiger) error {
	if conf == nil {
		return errors.New("conf is empty")
	}
	for _, confItem := range conf.AuthirizedApps {
		if confItem.Name == "" {
			return errors.New("app.name is empty")
		}
		if confItem.Secret == "" {
			return errors.New("app.secret is empty")
		}
		if confItem.AppID == "" {
			return errors.New("app.app_id is empty")
		}

		//validator
		errs := public.ValidatorHandler.Struct(confItem)
		if errs != nil {
			return errs
		}
	}
	return nil
}

func (s *SysConfigManage) getDBAPPConf(ischeck bool) (*dao.APPConfiger, error) {
	defer func() {
		if err := recover(); err != nil {
			public.SysLogger.Error("GetDBAPPConf_recover:%v", err)
		}
	}()

	apps, err := (&dao.GatewayAPP{}).GetAll()
	if err != nil {
		return nil, err
	}
	appConf := &dao.APPConfiger{AuthirizedApps: apps}
	if ischeck {
		if err := s.checkAppConf(appConf); err != nil {
			return nil, err
		}
	}
	return appConf, nil
}

func (s *SysConfigManage) writeFileAppConf(confPath string, appConf *dao.APPConfiger) error {
	var buffer bytes.Buffer
	e := toml.NewEncoder(&buffer)
	if err := e.Encode(appConf); err != nil {
		return err
	}
	bakPath := strings.Replace(confPath, ".toml", "_bak.toml", -1)
	os.Remove(bakPath)
	if err := os.Rename(confPath, bakPath); err != nil {
		return err
	}
	if err := ioutil.WriteFile(confPath, buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

//配置模块负载信息到ModuleRRMap+已运行模块周期刷新，直到IpContext停止
func (s *SysConfigManage) configModuleRR() error {
	modules := s.GetModuleConfig()
	for _, modulePointer := range modules.Module {
		currentModule := modulePointer
		go func(currentModule *dao.GatewayModule) {
			defer func() {
				if err := recover(); err != nil {
					public.SysLogger.Error("ConfigModuleRR_recover:%v", err)
				}
			}()
			if currentModule.Base.LoadType != "http" {
				return
			}
			t1 := time.NewTimer(0)
			ipList := []string{}
			ipWeightMap := map[string]int64{}
		Loop:
			for {
				select {
				case <-t1.C:
					newIPList := s.GetAvaliableIPList(currentModule.Base.Name)
					newIPWeightMap := s.GetConfIPWeightMap(currentModule, public.IPDefaultWeight)
					if !reflect.DeepEqual(ipList, newIPList) || !reflect.DeepEqual(ipWeightMap, newIPWeightMap) {
						Rw := public.NewWeightedRR(public.RRNginx)
						for _, ipAddr := range newIPList {
							w, ok := newIPWeightMap[ipAddr]
							if ok {
								Rw.Add(ipAddr, int(w))
							} else {
								Rw.Add(ipAddr, public.IPDefaultWeight)
							}
						}
						s.moduleRRMapLocker.Lock()
						s.moduleRRMap[currentModule.Base.Name] = Rw
						s.moduleRRMapLocker.Unlock()
					}
					ipList = newIPList
					ipWeightMap = newIPWeightMap
					t1.Reset(time.Millisecond * time.Duration(currentModule.LoadBalance.CheckInterval))
				case <-s.ConfigChangeNotice():
					t1.Stop()
					break Loop
				}
			}
		}(currentModule)
	}
	return nil
}

func (s *SysConfigManage) configModuleProxyMap() error {
	modules := s.GetModuleConfig()
	for _, modulePointer := range modules.Module {
		currentModule := modulePointer
		proxyFunc := func(rr public.RR) *httputil.ReverseProxy {
			proxy := &httputil.ReverseProxy{
				Director: func(req *http.Request) {
					if rHost, ok := rr.Next().(string); ok {
						req.URL.Scheme = "http"
						if req.TLS != nil {
							req.URL.Scheme = "https"
						}
						req.URL.Host = rHost
						req.Host = lib.GetStringConf("base.http.req_host")
					}
				},
				ModifyResponse: func(response *http.Response) error {
					var payload []byte
					var readErr error
					if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
						gr, err := gzip.NewReader(response.Body)
						if err != nil {
							public.ContextWarning(response.Request.Context(), lib.DLTagUndefind, map[string]interface{}{
								"func": "gzip.NewReader_error",
								"err":  err,
							})
						}
						payload, readErr = ioutil.ReadAll(gr)
						response.Header.Del(public.ContentEncoding)
					} else {
						payload, readErr = ioutil.ReadAll(response.Body)
					}
					if readErr != nil {
						public.ContextWarning(response.Request.Context(), lib.DLTagUndefind, map[string]interface{}{
							"func": "ioutil.ReadAll_error",
							"err":  readErr,
						})
						return readErr
					}

					//过滤前打点
					newPayload := payload
					public.ContextNotice(response.Request.Context(), lib.DLTagUndefind, map[string]interface{}{
						"func":       "modify_response_before",
						"payload":    public.Substr(string(payload), 0, 1<<8),
						"modulename": currentModule.Base.Name,
						"url":        response.Request.URL.String(),
						"header":     response.Header,
					})

					//过滤请求数据
					response.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
					response.ContentLength = int64(len(newPayload))
					response.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayload)), 10))
					if err := ModifyResponse(currentModule,response.Request, response);err != nil {
						return err
					}

					//过滤后打点
					public.ContextNotice(response.Request.Context(), lib.DLTagUndefind, map[string]interface{}{
						"func":       "modify_response_after",
						"newPayload": public.Substr(string(payload), 0, 1<<8),
						"modulename": currentModule.Base.Name,
						"url":        response.Request.URL.String(),
						"header":     response.Header,
					})
					return nil
				},
				Transport: &http.Transport{
					//请求下游的时间
					DialContext: (&net.Dialer{
						//限制建立TCP连接的时间
						Timeout: time.Duration(currentModule.LoadBalance.ProxyConnectTimeout) * time.Millisecond,
					}).DialContext,
					//最大空闲链接数
					MaxIdleConns: currentModule.LoadBalance.MaxIdleConn,
					//链接最大空闲时间
					IdleConnTimeout: time.Duration(currentModule.LoadBalance.IdleConnTimeout) * time.Millisecond,
					//限制读取response header的时间
					ResponseHeaderTimeout: time.Duration(currentModule.LoadBalance.ProxyHeaderTimeout) * time.Millisecond,
					//限制读取response body的时间
					ExpectContinueTimeout: time.Duration(currentModule.LoadBalance.ProxyBodyTimeout) * time.Millisecond,
				},
				ErrorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
					if err.Error() == "context canceled" {
						public.ContextWarning(req.Context(), lib.DLTagUndefind, map[string]interface{}{
							"func": "module_proxy_errorhandler",
							"err":  err,
							"url":  req.URL.String(),
						})
					} else {
						public.ContextError(req.Context(), lib.DLTagUndefind, map[string]interface{}{
							"func": "module_proxy_errorhandler",
							"err":  err,
							"url":  req.URL.String(),
						})
					}
					public.HTTPError(http.StatusGatewayTimeout, fmt.Sprint(err), w, req)
					return
				},
			}
			return proxy
		}
		s.moduleProxyFuncMapLocker.Lock()
		s.moduleProxyFuncMap[currentModule.Base.Name] = proxyFunc
		s.moduleProxyFuncMapLocker.Unlock()
	}
	return nil
}
