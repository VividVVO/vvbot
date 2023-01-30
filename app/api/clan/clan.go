package clan

import (
	"github.com/vivid-vvo/vvbot/app/service/clan"
	"github.com/vivid-vvo/vvbot/library/response"

	"github.com/gogf/gf/net/ghttp"
)

// API管理对象
type Controller struct{}

// 创建公会组请求参数
type ClanGroupCreateRequest struct {
	clan.ClanGroupCreateInput
}

// 获取用户公会列表请求参数
type GetUserClanListRequest struct {
	Qqid int64
}

// @summary 获取用户公会列表
// @tags    公会服务
// @produce json
// @param   qqid  formData string  true "qq号"
// @router  /pcr/getuserclanlist [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetUserClanList(r *ghttp.Request) {
	var data *GetUserClanListRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	userClanList, err := clan.GetUserClanList(r, data.Qqid)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", userClanList)
	}
}

// @summary 获取所有公会
// @tags    公会服务
// @produce json
// @router  /pcr/getallclan [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetAllClan(r *ghttp.Request) {
	userClanList, err := clan.GetAllClan(r)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", userClanList)
	}
}

// 通过公会ID获取群员数据请求参数
type GetClanGroupMembersRequest struct {
	ClanGroupID int `v:"min:1#公会ID不能为空"`
}

// @summary 通过公会ID获取群员数据
// @tags    公会服务
// @produce json
// @param   clanGroupID  formData string  true "公会id"
// @router  /pcr/getclangroupmembers [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetClanGroupMembers(r *ghttp.Request) {
	var data *GetClanGroupMembersRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	groupMembers, err := clan.GetClanGroupMembers(r, data.ClanGroupID)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", groupMembers)
	}
}

type ChangeUserDataRequest struct {
	QQid        int64 `v:"min:1#QQ号不能为空"`
	ClanGroupID int   `v:"min:1#公会ID不能为空"`
	GameName    string
	Role        int
}

// @summary 修改用户信息
// @tags    用户服务
// @produce json
// @param   qqid query string true "QQ号"
// @param   clangroupid query string true "公会ID"
// @param   gamename query string true "游戏名"
// @param   role query string true "权限"
// @router  /prc/changemembersdata [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) ChangeMembersData(r *ghttp.Request) {
	var data *ChangeUserDataRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := clan.ChangeMembersData(r, data.QQid, data.ClanGroupID, data.GameName, data.Role)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}

type MemberExitGroupRequest struct {
	QQid        int64
	ClanGroupID int `v:"min:1#公会ID不能为空"`
}

// @summary 成员退出公会
// @tags    用户服务
// @produce json
// @param   qqid query string true "QQ号"
// @param   clangroupid query string true "公会ID"
// @param   role query string true "权限"
// @router  /prc/memberexitgroup [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) MemberExitGroup(r *ghttp.Request) {
	var data *MemberExitGroupRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := clan.MemberExitGroup(r, data.QQid, data.ClanGroupID)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}

type ChangeClanInfoRequest struct {
	clan.ChangeClanInfoInput
}

// @summary 修改公会信息
// @tags    用户服务
// @produce json
// @param   groupid     int    true "公会ID"
// @param   groupname   string true "公会名称"
// @param   gameserver  string true "游戏服务器"
// @param   bindqqgroup int64 true "绑定Q群"
// @param   apikey      string "Apikey"
// @router  /prc/changeclaninfo [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) ChangeClanInfo(r *ghttp.Request) {
	var data *ChangeClanInfoRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := clan.ChangeClanInfo(r, data.ChangeClanInfoInput)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}

type DelClanGroupRequest struct {
	GroupId int `v:"min:1#公会ID不能为空"` //
}

// @summary 删除公会
// @tags    用户服务
// @produce json
// @param   groupid     int    true "公会ID"
// @router  /prc/delclangroup [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) DelClanGroup(r *ghttp.Request) {
	var data *DelClanGroupRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := clan.DelClanGroup(r, data.GroupId)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "")
}
