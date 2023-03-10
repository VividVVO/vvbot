package user

import (
	"github.com/vivid-vvo/vvbot/app/service/user"
	"github.com/vivid-vvo/vvbot/library/response"

	"github.com/gogf/gf/net/ghttp"
)

// 用户API管理对象
type Controller struct{}

// 注册请求参数，用于前后端交互参数格式约定
type SignUpRequest struct {
	user.SignUpInput
}

// @summary 用户注册接口
// @tags    用户服务
// @produce json
// @param   qqid  formData string  true "qq号"
// @param   password  formData string  true "用户密码"
// @param   nickname  formData string false "用户昵称"
// @router  /user/signup [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) SignUp(r *ghttp.Request) {
	var data *SignUpRequest
	// 这里没有使用Parse而是仅用GetStruct获取对象，
	// 数据校验交给后续的service层统一处理
	if err := r.GetStruct(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := user.SignUp(r, &data.SignUpInput); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

// 登录请求参数，用于前后端交互参数格式约定
type LoginRequest struct {
	QQid     string `v:"required#QQ号不能为空"`
	Password string
	Key      string
}

// @summary 用户登录接口
// @tags    用户服务
// @produce json
// @param   qqid formData string true "QQ号"
// @param   password formData string true "用户密码"
// @router  /user/signin [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) Login(r *ghttp.Request) {
	var data *LoginRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if data, err := user.Login(r, data.QQid, data.Password, data.Key); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", data)
	}
}

// @summary 判断用户是否已经登录
// @tags    用户服务
// @produce json
// @router  /user/issignedin [GET]
// @success 200 {object} response.JsonResponse "执行结果:`true/false`"
func (c *Controller) IsSignedIn(r *ghttp.Request) {
	response.JsonExit(r, 0, "", user.IsSignedIn(r))
}

// @summary 用户注销/退出接口
// @tags    用户服务
// @produce json
// @router  /user/signout [GET]
// @success 200 {object} response.JsonResponse "执行结果, 1: 未登录"
func (c *Controller) SignOut(r *ghttp.Request) {
	if err := user.SignOut(r); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "ok")
}

// 账号唯一性检测请求参数，用于前后端交互参数格式约定
type CheckQQidRequest struct {
	QQid string
}

// @summary 检测用户账号接口(唯一性校验)
// @tags    用户服务
// @produce json
// @param   qqid query string true "用户账号"
// @router  /user/checkqqid [GET]
// @success 200 {object} response.JsonResponse "执行结果:`true/false`"
func (c *Controller) CheckQQid(r *ghttp.Request) {
	var data *CheckQQidRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if data.QQid != "" && !user.CheckQQid(data.QQid) {
		response.JsonExit(r, 0, "账号已经存在", false)
	}
	response.JsonExit(r, 0, "", true)
}

// @summary 获取用户详情信息
// @tags    用户服务
// @produce json
// @router  /user/profile [GET]
// @success 200 {object} user.Entity "用户信息"
func (c *Controller) Profile(r *ghttp.Request) {
	data, err := user.GetProfile(r, 0)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "", data)
}

// 修改密码请求参数
type ChangePasswordRequest struct {
	user.ChangePasswordInput
}

// @summary 修改密码请求接口
// @tags    用户服务
// @produce json
// @param   qqid query string true "用户账号"
// @router  /user/changepassword [GET]
// @success 200 {object} response.JsonResponse "执行结果, 1: 修改失败"
func (c *Controller) ChangePassword(r *ghttp.Request) {
	var data *ChangePasswordRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := user.ChangePassword(r, &data.ChangePasswordInput)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}

// @summary 获取用户列表
// @tags    用户服务
// @produce json
// @router  /user/getuserlist [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetUserList(r *ghttp.Request) {
	userList, err := user.GetUserList()
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "", userList)
}

// 认证登录请求参数，用于前后端交互参数格式约定
type ChangeUserDataRequest struct {
	QQid     int64 `v:"min:1#QQ号不能为空"`
	NickName string
	Auth     int
}

// @summary 修改用户信息
// @tags    用户服务
// @produce json
// @param   qqid query string true "QQ号"
// @param   nickname query string true "用户别名"
// @param   auth query string true "权限"
// @router  /user/changeuserdata [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) ChangeUserData(r *ghttp.Request) {
	var data *ChangeUserDataRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := user.ChangeUserData(r, data.QQid, data.NickName, data.Auth)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}

// 请求参数，用于前后端交互参数格式约定
type DelUserRequest struct {
	QQid int64 `v:"min:1#QQ号不能为空"`
}

// @summary 删除用户
// @tags    用户服务
// @produce json
// @param   qqid query string true "QQ号"
// @router  /user/deltser [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) DelUser(r *ghttp.Request) {
	var data *DelUserRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := user.DelUser(r, data.QQid)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}
