package gvg

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/library/response"
)

// API管理对象
type Controller struct{}

// 报刀请求参数，用于前后端交互参数格式约定
type ReportCauseHarmRequest struct {
	gvg.ReportCauseHarmInput
}

// @summary 报刀
// @tags    公会战服务
// @produce json
// @param   qqid  formData string  true "QQ号"
// @param   damage  formData string  true "伤害"

// @router  /api/pcr/reportcauseharm [GET/POST]
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

type GetClanGvgRequest struct {
	ClanGroupID int `v:"min:1#公会ID不能为空"`
}

// @summary 获取公会战信息
// @tags    公会服务
// @produce json
// @param   clanGroupID  formData string  true "公会id"
// @router  /pcr/getclangvg [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetClanGvg(r *ghttp.Request) {
	var data *GetClanGvgRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	clanGvg, err := gvg.GetClanGvg(r, data.ClanGroupID)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", clanGvg)
	}
}

type GetChallengeAtQQRequest struct {
	QQID        int64
	ClanGroupID int `v:"min:1#公会ID不能为空"`
	TimeType    string
}

// @summary 获取用户战斗数据
// @tags    公会战服务
// @produce json
// @param   qqid  formData string  true "QQ号"
// @param   ClanGroupID  formData string  true "公会id"
// @param   TimeType  formData string  true "日期类型"
// @router  /api/pcr/getchallengeatqq [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetChallengeAtQQ(r *ghttp.Request) {
	var data *GetChallengeAtQQRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if challenges, err := gvg.GetChallengeAtQQ(r, data.QQID, data.ClanGroupID, data.TimeType); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", challenges)
	}
}

type GetAllChallengeRequest struct {
	ClanGroupID int `v:"min:1#公会ID不能为空"`
	TimeType    string
	StartTime   string
	EndTime     string
}

// @summary 获取所有战斗数据
// @tags    公会战服务
// @produce json
// @param   ClanGroupID  formData string  true "公会id"
// @param   TimeType  formData string  true "日期类型"
// @router  /api/pcr/getallchallenge [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetAllChallenge(r *ghttp.Request) {
	var data *GetAllChallengeRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if challenges, err := gvg.GetAllChallenge(r, data.ClanGroupID, data.TimeType, data.StartTime, data.EndTime); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", challenges)
	}
}

// @summary 获取所有战斗数据
// @tags    公会战服务
// @produce json
// @param   ClanGroupID  formData string  true "公会id"
// @param   TimeType  formData string  true "日期类型"
// @router  /api/pcr/DownAllChallengeToExcel [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) DownAllChallengeToExcel(r *ghttp.Request) {
	var data *GetAllChallengeRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if challengeToExcel, err := gvg.DownAllChallengeToExcel(r, data.ClanGroupID, data.TimeType); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		r.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", challengeToExcel.TimeStr))
		r.Response.Header().Set("filename", fmt.Sprintf("%s.xlsx", challengeToExcel.TimeStr))
		r.Response.Header().Set("Access-Control-Expose-Headers", "filename")
		r.Response.WriteExit(challengeToExcel.ExcelBuffer)
		// response.JsonExit(r, 0, "ok", challenges)
	}
}

type RemindChallengeRequest struct {
	ClanGroupID int     `v:"min:1#公会ID不能为空"`
	QQIDlist    []int64 `v:"required#成员不能为空"`
	Type        int
}

// @summary 提醒战斗
// @tags    公会战服务
// @produce json
// @param   ClanGroupID  formData string  true "公会id"
// @param   QQIDlist  formData string  true "QQ号列表"
// @router  /api/pcr/remindchallenge [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) RemindChallenge(r *ghttp.Request) {
	var data *RemindChallengeRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gvg.RemindChallenge(r, data.ClanGroupID, data.QQIDlist, data.Type); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

type ChangeUserChallengeRequest struct {
	gvg.ChangeUserChallengeInput
}

// @summary 修改用户战斗数据
// @tags    公会战服务
// @produce json
// @param   GroupId formData int true "公会ID"
// @param   ChallengeId formData int true "战斗ID"
// @param   ChallengeDamage formData true "战斗伤害"
// @param   BossCycle formData string true "周目"
// @param   BossNum formData string true "BOSS号"
// @param   Meassage formData string true "留言"
// @router  /api/pcr/changeuserchallenge [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) ChangeUserChallenge(r *ghttp.Request) {
	var data *ChangeUserChallengeRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gvg.ChangeUserChallenge(r, data.ChangeUserChallengeInput); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

type DelUserChallengeRequest struct {
	ChallengeId int `v:"min:1#战斗ID不能为空"`
	GroupId     int `v:"min:1#公会ID不能为空"`
}

// @summary 删除用户战斗数据
// @tags    公会战服务
// @produce json
// @param   GroupId formData int true "公会ID"
// @param   ChallengeId formData int true "战斗ID"
// @router  /api/pcr/deluserchallenge [GET/POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) DelUserChallenge(r *ghttp.Request) {
	var data *ChangeUserChallengeRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gvg.DelUserChallenge(r, data.GroupId, data.ChallengeId); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok")
	}
}

type GetAllSlStateRequest struct {
	ClanGroupID int `v:"min:1#公会ID不能为空"`
	TimeType    string
}

// @summary 获取所有成员sl状态
// @tags    用户服务
// @produce json
// @param   groupid     int    true "公会ID"
// @param   timetype  formData string  true "日期类型"
// @router  /prc/getallslstate [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) GetAllSlState(r *ghttp.Request) {
	var data *GetAllSlStateRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}

	if clanGroupData, err := gvg.GetAllSlState(r, data.ClanGroupID, data.TimeType); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "ok", clanGroupData)
	}
}
