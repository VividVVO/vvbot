package gvg

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/library/response"
)

// API管理对象
type Controller struct{}

// 创建公会请求参数，用于前后端交互参数格式约定
type GvgGroupCreateRequest struct {
	gvg.GvgGroupCreateInput
}

// @summary 创建一个新的公会战
// @tags    公会战服务
// @produce json
// @param   groupname  formData string  true "公会组ID"
// @param   gameserver  formData string  true "游戏服务器"
// @param   gvgname  formData string  true "公会战名称"
// @param   gvgstarttime  formData string  true "公会战开始时间"
// @param   gvgendtime  formData string  true "公会战结束时间"
// @router  /user/gvggroupcreate [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GvgGroupCreate(r *ghttp.Request) {
	//noinspection GoUnresolvedReference
	var data *GvgGroupCreateRequest
	if err := r.GetStruct(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gvg.GvgGroupCreate(r, &data.GvgGroupCreateInput, 0); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

// 报刀请求参数，用于前后端交互参数格式约定
type ReportCauseHarmRequest struct {
	gvg.ReportCauseHarmInput
}

// @summary 报刀
// @tags    公会战服务
// @produce json
// @param   qqid  formData string  true "QQ号"
// @param   damage  formData string  true "伤害"

// @router  /user/reportcauseharm [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) ReportCauseHarm(r *ghttp.Request) {
	//noinspection GoUnresolvedReference
	var data *ReportCauseHarmRequest
	if err := r.GetStruct(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gvg.ReportCauseHarm(r, &data.ReportCauseHarmInput); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}
