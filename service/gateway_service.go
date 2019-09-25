package service

import (
	"fmt"
	"github.com/didi/gatekeeper/dao"
	"github.com/didi/gatekeeper/public"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//dltag常量
const (
	DLTagMatchRuleFailure     = " match_rule_failure"
	DLTagMatchRuleSuccess     = " match_rule_success"
	DLTagAccessControlFailure = " access_control_failure"
	DLTagAccessControlUndef   = " access_control_undef"
	DLTagAccessControlSuccess = " access_control_success"
	DLTagLoadBalanceFailure   = " load_balance_failure"
	DLTagLoadBalanceSuccess   = " load_balance_suceess"
	DLTagSsoHandlerFailure    = " sso_handler_failure"
	DLTagSsoHandlerSuccess    = " sso_handler_success"

	HeaderKeyUserPerm     = "didi-header-userperm"
	HeaderKeyUserGroup    = "didi-header-usergroup"
	HeaderKeyUserGroupKey = "didi-header-usergroupkey"
	UserPermCtxKey        = "user_perm_key"
	HeaderKeyUserCityPerm = "didi-header-usercityperm"
)

//GateWayService 网关核心服务
type GateWayService struct {
	currentModule *dao.GatewayModule
	w             http.ResponseWriter
	req           *http.Request
}

//NewGateWayService 构建一个服务
func NewGateWayService(w http.ResponseWriter, req *http.Request) *GateWayService {
	return &GateWayService{
		w:   w,
		req: req,
	}
}

//CurrentModule 当前模块
func (o *GateWayService) CurrentModule() *dao.GatewayModule {
	return o.currentModule
}

//SetCurrentModule 设置当前模块
func (o *GateWayService) SetCurrentModule(currentModule *dao.GatewayModule) {
	o.currentModule = currentModule
}


//AccessControl 权限验证
func (o *GateWayService) AccessControl() error {
	if o.currentModule.AccessControl == nil {
		return nil
	}
	ctx := public.NewContext(o.w, o.req)
	var errmsg string
	switch {
	case !AuthModuleOpened(o, ctx):
		public.ContextNotice(o.req.Context(), DLTagAccessControlSuccess, map[string]interface{}{
			"msg": "access_control_not_open",
		})
		return nil
	case AuthInBlackIPList(o, ctx):
		public.ContextNotice(o.req.Context(), DLTagAccessControlFailure, map[string]interface{}{
			"msg": "AuthInBlackIPList",
		})
		return errors.New("msg:AuthInBlackIPList")
	case AuthInWhiteIPList(o, ctx):
		public.ContextNotice(o.req.Context(), DLTagAccessControlSuccess, map[string]interface{}{
			"msg": "AuthWhiteIPList_success",
		})
		return nil
	case AuthInWhiteHostList(o, ctx):
		public.ContextNotice(o.req.Context(), DLTagAccessControlSuccess, map[string]interface{}{
			"msg": "AuthWhiteHostList_success",
		})
		return nil
	case AuthRegisterFunc(o, &errmsg):
		public.ContextNotice(o.req.Context(), DLTagAccessControlSuccess, map[string]interface{}{
			"msg": "AuthRegisterFunc_success",
		})
		return nil
	}
	if errmsg==""{
		errmsg="auth_failure"
	}
	public.ContextWarning(o.req.Context(), DLTagAccessControlFailure, map[string]interface{}{
		"msg": errmsg,
	})
	return errors.New(errmsg)
}

//AuthModuleOpened 是否启用模块权限校验？
func AuthModuleOpened(o *GateWayService, ctx *public.Context) bool {
	if o.currentModule.AccessControl.Open == 1 {
		return true
	}
	return false
}

//AuthAppToken app的签名校验
func AuthAppToken(m *dao.GatewayModule, req *http.Request) (bool,error) {
	ctx:=public.NewContext(nil,req)
	if err:=AuthAppSign(ctx);err!=nil {
		return false,err
	}
	if err := AfterAuthLimit(ctx); err != nil {
		return false,err
	}
	return true,nil
}

//AuthInBlackIPList 黑名单验证
func AuthInBlackIPList(o *GateWayService, ctx *public.Context) bool {
	clientIP := public.RemoteIP(o.req)
	blackList := strings.Split(o.currentModule.AccessControl.BlackList, ",")
	if public.AuthIPList(clientIP, blackList) {
		public.ContextNotice(o.req.Context(), DLTagAccessControlUndef, map[string]interface{}{
			"msg":       "AuthInBlackIPList",
			"clientIP":  clientIP,
			"blackList": blackList,
		})
		return true
	}
	return false
}

//AuthInWhiteIPList 验证白名单
func AuthInWhiteIPList(o *GateWayService, ctx *public.Context) bool {
	clientIP := public.RemoteIP(o.req)
	whiteList := strings.Split(o.currentModule.AccessControl.WhiteList, ",")
	if public.AuthIPList(clientIP, whiteList) {
		return true
	}
	return false
}

//AuthInWhiteHostList 验证host白名单
func AuthInWhiteHostList(o *GateWayService, ctx *public.Context) bool {
	hostname, err := os.Hostname()
	if err != nil {
		return false
	}
	whiteHostname := strings.Split(o.currentModule.AccessControl.WhiteHostName, ",")
	if public.AuthIPList(hostname, whiteHostname) {
		return true
	}
	return false
}

//AuthRegisterFunc 验证注册函数
func AuthRegisterFunc(o *GateWayService, errmsg *string) bool {
	for _, rf := range BeforeRequestAuthRegisterFuncs {
		flag,err:=rf(o.currentModule, o.req)
		if flag {
			return true
		}
		*(errmsg)=err.Error()
	}
	return false
}

//AfterAuthLimit 验证后限流
func AfterAuthLimit(authCtx *public.Context) error {
	appID := authCtx.Query("app_id")
	return AppAuth(appID, authCtx)
}

//AppAuth app验证
func AppAuth(appID string, authCtx *public.Context) error {
	appConfig, err := SysConfMgr.GetAppConfigByAPPID(appID)
	if err != nil {
		return err
	}
	v:=authCtx.Req.Context().Value(public.ContextKey("request_url"))
	reqPath,ok := v.(string)
	if !ok{
		reqPath = ""
	}
	if !public.InOrPrefixStringList(reqPath, strings.Split(appConfig.OpenAPI, ",")) {
		errmsg := "You don't have rights for this path:" + reqPath + " - " + appConfig.OpenAPI
		return errors.New(errmsg)
	}

	//限速器
	limiter := public.FlowLimiterHandler.GetAPPLimitVisitor(appID, appConfig.QPS)
	if appConfig.QPS > 0 && limiter.Allow() == false {
		errmsg := fmt.Sprintf("QPS limit : %d, %d", int64(limiter.Limit()), limiter.Burst())
		return errors.New(errmsg)
	}

	if appConfig.GroupID > 0 {
		authCtx.Req.Header.Add(HeaderKeyUserGroup, strconv.Itoa(int(appConfig.GroupID)))
		authCtx.Req.Header.Add(HeaderKeyUserGroupKey, public.UserGroupPerfix+strconv.Itoa(int(appConfig.GroupID)))
	}

	counter := public.FlowCounterHandler.GetAPPCounter(appID)
	if appConfig.TotalQueryDaily > 0 && counter.TotalCount > appConfig.TotalQueryDaily {
		errmsg := fmt.Sprintf("total query daily limit: %d", appConfig.TotalQueryDaily)
		return errors.New(errmsg)
	}
	counter.Increase(authCtx.Req.Context())
	return nil
}

//AuthAppSign 验证app签名
func AuthAppSign(c *public.Context) error {
	clientSign := c.Query("sign")
	appID := c.Query("app_id")
	if appID == "" {
		return errors.New(fmt.Sprintf("AuthAppSign -error:%v",
				"app_id empty"))
	}
	appConfig, err := SysConfMgr.GetAppConfigByAPPID(appID)
	if err != nil {
		return errors.New(fmt.Sprintf(
				"AuthAppSign -error:%v -app_id:%v -sign:%v",
				"GetAppConfigByAPPID error", appID,clientSign, ))
	}
	if appConfig.Secret == "" {
		return errors.New(fmt.Sprintf(
			"AuthAppSign -error:%v -app_id:%v -sign:%v",
			"Secret empty", appID,clientSign, ))
	}

	if appConfig.WhiteIps != "" &&
		public.AuthIPList(public.RemoteIP(c.Req), strings.Split(appConfig.WhiteIps, ",")) {
		return nil
	}

	//todo 自定义sign生成规则
	signKey := appConfig.Secret
	if signKey != clientSign {
		return errors.New(fmt.Sprintf(
			"AuthAppSign -error:%v -app_id:%v -sign:%v",
			"sign error", appID,clientSign, ))
	}
	return nil
}

//LoadBalance 请求负载
func (o *GateWayService) LoadBalance() (*httputil.ReverseProxy, error) {
	ipList, err := SysConfMgr.GetModuleIPList(o.currentModule.Base.Name)
	if err != nil {
		public.ContextWarning(o.req.Context(), DLTagLoadBalanceFailure, map[string]interface{}{
			"msg":             err,
			"modulename":      o.currentModule.Base.Name,
			"availableIpList": SysConfMgr.GetAvaliableIPList(o.currentModule.Base.Name),
		})
		return nil, errors.New("get_iplist_error")
	}
	if len(ipList) == 0 {
		public.ContextWarning(o.req.Context(), DLTagLoadBalanceFailure, map[string]interface{}{
			"msg":             "empty_iplist_error",
			"modulename":      o.currentModule.Base.Name,
			"availableIpList": SysConfMgr.GetAvaliableIPList(o.currentModule.Base.Name),
		})
		return nil, errors.New("empty_iplist_error")
	}
	proxy, err := o.GetModuleHTTPProxy()
	if err != nil {
		public.ContextWarning(o.req.Context(), DLTagLoadBalanceFailure, map[string]interface{}{
			"msg":       err,
			"module":    o.currentModule.Base.Name,
		})
		return nil, err
	}
	return proxy, nil
}

//GetModuleHTTPProxy 获取模块的代理
func (o *GateWayService) GetModuleHTTPProxy() (*httputil.ReverseProxy, error) {
	proxy,err:=SysConfMgr.GetModuleHTTPProxy(o.currentModule.Base.Name)
	if err != nil {
		public.ContextWarning(o.req.Context(), DLTagLoadBalanceFailure, map[string]interface{}{
			"err":       err,
			"module":    o.currentModule.Base.Name,
		})
		return &httputil.ReverseProxy{}, err
	}
	return proxy,nil
}

//MatchRule 匹配规则
func (o *GateWayService) MatchRule() error {
	var currentModule *dao.GatewayModule
	modules := SysConfMgr.GetModuleConfig()
Loop:
	for _, module := range modules.Module {
		if module.Base.LoadType != "http" {
			continue
		}
		for _, matchRule := range module.MatchRule {
			urlStr := o.req.URL.Path
			if matchRule.Type == "url_prefix" && strings.HasPrefix(urlStr, matchRule.Rule+"/") {
				currentModule = module
				//提前检测，减少资源消耗
				if matchRule.URLRewrite ==""{
					break Loop
				}
				for _,uw  := range strings.Split(matchRule.URLRewrite, ",") {
					uws := strings.Split(uw, " ")
					if len(uws) == 2 {
						re, regerr := regexp.Compile(uws[0])
						if regerr!=nil{
							return regerr
						}
						rep := re.ReplaceAllString(urlStr, uws[1])
						o.req.URL.Path = rep
						public.ContextNotice(o.req.Context(), DLTagMatchRuleSuccess, map[string]interface{}{
							"url":       o.req.RequestURI,
							"write_url": rep,
						})
						if o.req.URL.Path!=urlStr{
							break
						}
					}
				}
				break Loop
			}
		}
	}
	if currentModule == nil {
		public.ContextWarning(o.req.Context(), DLTagMatchRuleFailure, map[string]interface{}{
			"msg": "module_not_found",
			"url": o.req.RequestURI,
		})
		return errors.New("module not found")
	}
	public.ContextNotice(o.req.Context(), DLTagMatchRuleSuccess, map[string]interface{}{
		"url": o.req.RequestURI,
	})
	o.SetCurrentModule(currentModule)
	return nil
}