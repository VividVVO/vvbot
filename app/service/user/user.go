package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/model/user_login"
	"github.com/vivid-vvo/vvbot/app/service/check"
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
	Qqid           int64  `json:"qqid"`            //   QQ号
	Nickname       string `json:"nickname"`        //   用户名
	AuthorityGroup int    `json:"auth"`            //   所在权限组
	ClanGroupId    int64  `json:"clan_group_id"`   //   公会组ID
	LastLoginTime  int64  `json:"last_login_time"` //   最后登录时间
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
	QQid, err := strconv.ParseInt(data.QQid, 10, 64)
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

type LoginData struct {
	Token    string   `json:"token"`
	Roles    []string `json:"roles"`
	QQID     int64    `json:"qqid"`
	NickName string   `json:"nickname"`
}

// 用户登录，成功返回用户信息，否则返回nil;
func Login(r *ghttp.Request, qqid string, password string, authCookie string) (*LoginData, error) {
	var one *user.Entity
	var err error
	if password != "" {
		one, err = user.FindOne("qqid=? and password=?", qqid, fmt.Sprintf("%x", sha256.Sum256([]byte(password))))
		if err != nil {
			return nil, errors.New("内部错误")
		}
		if one == nil {
			return nil, errors.New("账号或密码错误")
		}
	} else if authCookie != "" {
		one, err = user.FindOne("qqid=? and login_code=?", qqid, authCookie)
		if err != nil {
			return nil, errors.New("内部错误")
		}
		if one == nil {
			return nil, errors.New("key错误")
		}
	} else {
		return nil, errors.New("参数错误")
	}
	nowTimeS := time.Now().Unix()
	entity := new(user_login.Entity)
	entity.AuthCookie = guid.S()
	entity.LoginTime = nowTimeS
	entity.LoginIp = r.GetClientIp()
	entity.Qqid = one.Qqid
	if _, err := user_login.Insert(entity); err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	one.LastLoginTime = nowTimeS
	one.LastLoginIp = r.GetClientIp()
	if _, err := user.Update(one, "qqid", qqid); err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	r.Cookie.Set("authkey", entity.AuthCookie)
	var roles []string
	iQQID, _ := strconv.ParseInt(qqid, 10, 64)
	if check.CheckIsMaster(iQQID) {
		roles = append(roles, "superadmin")
	} else {
		switch one.AuthorityGroup {
		case check.AuthUser:
			roles = append(roles, "member")
		case check.AuthAdmin:
			roles = append(roles, "admin")
		case check.AuthSuperAdmin:
			roles = append(roles, "superadmin")
		}
	}
	loginData := &LoginData{
		Token:    entity.AuthCookie,
		Roles:    roles,
		QQID:     one.Qqid,
		NickName: one.Nickname,
	}
	return loginData, nil
}

// 修改用户公会组ID;
func ChangeClanGroupId(r *ghttp.Request, qqid int64, groupId int) error {
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
	qqid, err := strconv.ParseInt(data.QQid, 10, 64)
	if err != nil {
		return errors.New("内部错误")
	}
	err = user.ChangePassword(qqid, data.Password)
	if err != nil {
		return err
	}
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
func GetProfile(r *ghttp.Request, qqid int64) (*UserInfo, error) {
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

// 登录信息、缓存信息
func GetLoginData2(r *ghttp.Request) user_login.Entity {
	var loginData user_login.Entity
	r.Session.GetStruct(SessionUserLogin, &loginData)

	return loginData
}

// 获取用户列表
func GetUserList() ([]*user.Entity, error) {
	userList, err := user.GetUserList()
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// 修改用户信息
func ChangeUserData(r *ghttp.Request, qqid int64, nickName string, auth int) error {
	thisqqid := GetLoginData2(r).Qqid
	profile, err := GetProfile(r, qqid)
	if err != nil {
		return err
	}
	if profile.AuthorityGroup != auth {
		switch auth {
		case check.AuthUser:
		case check.AuthSuperAdmin:
		case check.AuthAdmin:
		case check.AuthClanAdmin:
		case check.AuthGvgAdmin:
		default:
			return errors.New("参数错误")
		}
		if !check.CheckIsMaster(thisqqid) {
			if profile.AuthorityGroup == check.AuthSuperAdmin {
				return errors.New("权限不足")
			}
		}
		if !check.CheckAuthorityGroup(thisqqid, check.AuthSuperAdmin, 0) {
			return errors.New("权限不足")
		}
	} else if thisqqid != qqid {
		// 必须拥有管理员权限才能修改用户信息
		if !check.CheckAuthorityGroup(thisqqid, check.AuthAdmin, 0) {
			return errors.New("权限不足")
		}
	}
	return user.ChangeUserData(qqid, nickName, auth)
}

// 修改用户信息
func DelUser(r *ghttp.Request, qqid int64) error {
	thisqqid := GetLoginData2(r).Qqid
	profile, err := GetProfile(r, qqid)
	if err != nil {
		return err
	}
	if !check.CheckIsMaster(thisqqid) {
		if profile.AuthorityGroup == check.AuthSuperAdmin {
			return errors.New("权限不足")
		} else if profile.AuthorityGroup >= check.AuthAdmin {
			if !check.CheckAuthorityGroup(thisqqid, check.AuthSuperAdmin, 0) {
				return errors.New("权限不足")
			}
		} else {
			if !check.CheckAuthorityGroup(thisqqid, check.AuthAdmin, 0) {
				return errors.New("权限不足")
			}
		}
	}
	tx, err := g.DB().Begin()
	if err != nil {
		return errors.New("内部错误")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Table("user").Delete("qqid", qqid)
	if err != nil {
		return errors.New("删除失败")
	}
	_, err = tx.Table("user_login").Delete("qqid", qqid)
	if err != nil {
		return errors.New("删除失败")
	}
	_, err = tx.Table("clan_member").Delete("qqid", qqid)
	if err != nil {
		return errors.New("删除失败")
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("删除失败")
	}
	return nil
}
