package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/model/user_login"
	"strconv"
	"time"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/util/guid"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

const (
	SessionUserLogin = "user_login"
)

// 注册输入参数
type SignUpInput struct {
	QQid     string `v:"required|length:3,13#账号不能为空|账号长度应当在:min到:max之间"`
	Password string `v:"required|length:6,16#请输入密码|密码长度应当在:min到:max之间"`
	Nickname string
}

type ChangePasswordInput struct {
	QQid     string `v:"required|length:3,13#账号不能为空|账号长度应当在:min到:max之间"`
	Password string `v:"required|length:6,16#请输入密码|密码长度应当在:min到:max之间"`
}

// 用户信息
type UserInfo struct {
	Qqid           int    `json:"qqid"`            //   QQ号
	Nickname       string `json:"nickname"`        //   用户名
	AuthorityGroup int    `json:"authority_group"` //   所在权限组
	ClanGroupId    int    `json:"clan_group_id"`   //   公会组ID
	LastLoginTime  int    `json:"last_login_time"` //   最后登录时间
	LastLoginIp    string `json:"last_login_ip"`   //   最后登录IP
	ClanGroupName  string `json:"clan_group_name"`
	GvgName        string `json:"gvg_name"`
}
type UserData struct {
	user.Entity
	AuthCookie string `orm:"auth_cookie"     json:"auth_cookie"` // 登录认证Cookie
}

// 用户注册
func SignUp(r *ghttp.Request, data *SignUpInput) error {
	// 输入参数检查
	if e := gvalid.CheckStruct(data, nil); e != nil {
		return errors.New(e.FirstString())
	}
	// 账号唯一性数据检查
	if !CheckQQid(data.QQid) {
		return errors.New(fmt.Sprintf("账号 %s 已经存在", data.QQid))
	}
	QQid, err := strconv.Atoi(data.QQid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	err = user.SignUp(QQid, data.Password, data.Nickname, r.GetClientIp())
	if err != nil {
		return err
	}
	return nil
}

// 判断用户是否已经登录
func IsSignedIn(r *ghttp.Request) bool {
	authCookie := r.Cookie.Get("authkey")
	one, err := user_login.FindOne("auth_cookie", authCookie)
	if err != nil {
		return false
	}
	if one == nil {
		return false
	}
	r.Session.Set(SessionUserLogin, one)
	return true
}

// 获取登录信息
func GetLoginData(r *ghttp.Request) (*user_login.Entity, error) {
	authCookie := r.Cookie.Get("authkey")
	one, err := user_login.FindOne("auth_cookie", authCookie)
	if err != nil {
		return one, errors.New("内部错误")
	}
	r.Session.Set(SessionUserLogin, one)
	return one, nil
}

// 用户登录，成功返回用户信息，否则返回nil;
func SignIn(r *ghttp.Request, qqid string, password string) error {
	one, err := user.FindOne("qqid=? and password=?", qqid, fmt.Sprintf("%x", sha256.Sum256([]byte(password))))
	if err != nil {
		return errors.New("内部错误")
	}
	if one == nil {
		return errors.New("账号或密码错误")
	}
	nowTimeS := time.Now().Unix()
	var userData UserData
	if err := gconv.Struct(one, &userData); err != nil {
		return err
	}
	entity := new(user_login.Entity)
	entity.AuthCookie = guid.S()
	entity.LoginTime = int(nowTimeS)
	entity.LoginIp = r.GetClientIp()
	entity.Qqid = one.Qqid
	userData.AuthCookie = entity.AuthCookie
	if _, err := user_login.Insert(entity); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	one.LastLoginTime = int(nowTimeS)
	one.LastLoginIp = r.GetClientIp()
	if _, err := user.Update(one, "qqid", qqid); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	r.Cookie.Set("authkey", entity.AuthCookie)
	return nil
}

// 修改用户公会组ID;
func ChangeClanGroupId(r *ghttp.Request, qqid int, groupId int) error {
	if qqid == 0 {
		qqid = GetLoginData2(r).Qqid
	}
	if err := user.ChangeClanGroupId(qqid, groupId); err != nil {
		return err
	}
	return nil
}

// 修改用户密码，成功返回用户信息，否则返回nil;
func ChangePassword(r *ghttp.Request, data *ChangePasswordInput) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return errors.New("内部错误")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			r.Cookie.Remove("authkey")
		}
	}()
	_, err = tx.Table("user").Update(g.Map{"password": fmt.Sprintf("%x", sha256.Sum256([]byte(data.Password)))}, "qqid", data.QQid)
	if err != nil {
		return errors.New("修改密码失败！")
	}
	_, err = tx.Table("user_login").Delete("qqid", data.QQid)

	return nil
}

// 用户注销
func SignOut(r *ghttp.Request) error {
	authCookie := r.Cookie.Get("authkey")
	_, err := user_login.Delete("auth_cookie", authCookie)
	if err != nil {
		return errors.New("内部错误")
	}
	r.Session.Remove(SessionUserLogin)
	return nil
}

// 检查账号是否符合规范(目前仅检查唯一性),存在返回false,否则true
func CheckQQid(qqid string) bool {
	if i, err := user.FindCount("qqid", qqid); err != nil {
		return false
	} else {
		return i == 0
	}
}

// 检查昵称是否符合规范(目前仅检查唯一性),存在返回false,否则true
func CheckNickName(nickname string) bool {
	if i, err := user.FindCount("nickname", nickname); err != nil {
		return false
	} else {
		return i == 0
	}
}

// 获得用户信息详情
func GetProfile(r *ghttp.Request, qqid int) (*UserInfo, error) {
	if qqid == 0 {
		qqid = GetLoginData2(r).Qqid
	}
	profile, err := user.GetProfile(qqid)
	if err != nil {
		return nil, errors.New("内部错误")
	}
	var userInfo *UserInfo
	if err := gconv.Struct(profile, &userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// 0 普通用户, 1 管理员
func GetUserAuthorityGroup(r *ghttp.Request, qqid int) (int, error) {
	if qqid == 0 {
		qqid = GetLoginData2(r).Qqid
	}
	user, err := GetProfile(r, qqid)
	if err != nil {
		return 0, err
	}
	var authorityGroup int
	if user.AuthorityGroup == 0 {
		authorityGroup = 0
	} else if user.AuthorityGroup == 100 {
		authorityGroup = 1
	}
	return authorityGroup, nil
}

// 登录信息、缓存信息
func GetLoginData2(r *ghttp.Request) user_login.Entity {
	var loginData user_login.Entity
	r.Session.GetStruct(SessionUserLogin, &loginData)

	return loginData
}
