package dashboard_controller

import (
	"encoding/base64"
	"github.com/didi/gatekeeper/dashboard_middleware"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type OAuthController struct{}

func OAuthRegister(group *gin.RouterGroup) {
	oauth := &OAuthController{}
	group.POST("/tokens", oauth.Tokens)
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oauth *OAuthController) Tokens(c *gin.Context) {
	params := &model.TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		dashboard_middleware.ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	//fmt.Println("appSecret", string(appSecret))

	//  取出 app_id secret
	//  生成 app_list
	//  匹配 app_id
	//  基于 jwt生成token
	//  生成 output
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		dashboard_middleware.ResponseError(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}

	appList := handler.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(public.JwtExpires * time.Second).In(lib.TimeLocation).Unix(),
			}
			token,err:=public.JwtEncode(claims)
			if err != nil {
				dashboard_middleware.ResponseError(c, 2004, err)
				return
			}
			output := &model.TokensOutput{
				ExpiresIn:public.JwtExpires,
				TokenType:"Bearer",
				AccessToken:token,
				Scope:"read_write",
			}
			dashboard_middleware.ResponseSuccess(c, output)
			return
		}
	}
	dashboard_middleware.ResponseError(c, 2005,errors.New("未匹配正确APP信息"))
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminlogin *OAuthController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	dashboard_middleware.ResponseSuccess(c, "")
}
