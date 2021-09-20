package dashboard_controller

import (
	"encoding/json"
	"fmt"
	"github.com/didi/gatekeeper/dashboard_middleware"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminlogin *AdminController) AdminInfo(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &model.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	//1. 读取sessionKey对应json 转换为结构体
	//2. 取出数据然后封装输出结构体

	//Avatar       string    `json:"avatar"`
	//Introduction string    `json:"introduction"`
	//Roles        []string  `json:"roles"`
	out := &model.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	dashboard_middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminlogin *AdminController) ChangePwd(c *gin.Context) {
	params := &model.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	//1. session读取用户信息到结构体 sessInfo
	//2. sessInfo.ID 读取数据库信息 adminInfo
	//3. params.password+adminInfo.salt sha256 saltPassword
	//4. saltPassword==> adminInfo.password 执行数据保存

	//session读取用户信息到结构体
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &model.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	//从数据库中读取 adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &model.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, (&model.Admin{UserName: adminSessionInfo.UserName}))
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}

	//生成新密码 saltPassword
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword

	//执行数据保存
	if err := adminInfo.Save(c, tx); err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, "")
}
