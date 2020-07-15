package clan

import (
	"github.com/vivid-vvo/vvbot/app/service/clan"
	"github.com/vivid-vvo/vvbot/library/response"

	"github.com/gogf/gf/util/gvalid"

	"github.com/gogf/gf/net/ghttp"
)

// API管理对象
type Controller struct{}

// 创建公会组请求参数
type ClanGroupCreateRequest struct {
	clan.ClanGroupCreateInput
}

// @summary 创建公会
// @tags    公会服务
// @produce json
// @param   groupname  formData string  true "公会名"
// @param   gameserver  formData string  true "游戏服务器"
// @router  /user/clangroupcreate [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) ClanGroupCreate(r *ghttp.Request) {
	var data *ClanGroupCreateRequest
	// 这里没有使用Parse而是仅用GetStruct获取对象，
	// 数据校验交给后续的service层统一处理
	if err := r.GetStruct(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := clan.ClanGroupCreate(r, &data.ClanGroupCreateInput); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

// 加入公会组请求参数
type JoinClanRequest struct {
	GroupName string `v:"required#公会名不能为空"` //   公会ID
}

// @summary 加入公会
// @tags    公会服务
// @produce json
// @param   groupname  formData string  true "公会ID"
// @param   gamename  formData string  true "游戏名"
// @router  /user/joinclan [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) JoinClan(r *ghttp.Request) {
	var data *JoinClanRequest
	// 这里没有使用Parse而是仅用GetStruct获取对象，
	// 数据校验交给后续的service层统一处理
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := clan.JoinClan(r, 0, data.GroupName); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

// 加入公会组请求参数
type JoinClanToRequest struct {
	Qqid      int    `v:"min:1#QQ号不能为空"`
	GroupName string `v:"required#公会组名不能为空"`
}

// @summary 指定QQ加入公会
// @tags    公会服务
// @produce json
// @param   qqid  formData string  true "QQ号"
// @param   groupname  formData string  true "公会ID"
// @param   gamename  formData string  true "游戏名"
// @router  /user/joinclanto [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) JoinClanTo(r *ghttp.Request) {
	var data *JoinClanToRequest
	// 这里没有使用Parse而是仅用GetStruct获取对象，
	// 数据校验交给后续的service层统一处理
	if err := r.GetStruct(&data); err != nil {
		response.JsonExit(r, 1, err.Error())

	}
	if e := gvalid.CheckStruct(data, nil); e != nil {
		response.JsonExit(r, 1, e.FirstString())
	}
	if err := clan.JoinClan(r, data.Qqid, data.GroupName); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}
