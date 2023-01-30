package iotqq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"io/ioutil"
	"net/http"
	"net/url"
)

type QQ struct {
	Cont int
}

type MemberData struct {
	Count      int      `json:"Count"`
	LastUin    int      `json:"LastUin"`
	MemberList []Member `json:"MemberList"`
}

type GroupData struct {
	Count     int     `json:"Count"`
	NextToken string  `json:"NextToken"`
	GroupList []Group `json:"TroopList"`
}

type Group struct {
	GroupId          int64  `json:"GroupId"`
	GroupMemberCount int    `json:"GroupMemberCount"`
	GroupName        string `json:"GroupName"`
	GroupNotice      string `json:"GroupNotice"`
	GroupOwner       int64  `json:"GroupOwner"`
	GroupTotalCount  int    `json:"GroupTotalCount"`
}

type Member struct {
	JoinTime      int64  `json:"JoinTime"`
	AutoRemark    string `json:"AutoRemark"`
	GroupCard     string `json:"GroupCard"`
	LastSpeakTime int64  `json:"LastSpeakTime"`
	Qqid          int64  `json:"MemberUin"`
	NickName      string `json:"NickName"`
	Age           int    `json:"Age"`
	GroupAdmin    int    `json:"GroupAdmin"`
}

func getQQ() string {
	return qq
}

func sendPic(ToUser int64, SendToType int, Content string, PicBase64Buf string, PicUrl string) {
	//发送图文信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = ToUser
	tmp["sendToType"] = SendToType
	tmp["sendMsgType"] = "PicMsg"
	tmp["picBase64Buf"] = PicBase64Buf
	tmp["fileMd5"] = ""
	tmp["picUrl"] = PicUrl
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = 0
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
func send(ToUser int64, SendToType int, Content string, Qqid int64) {
	//发送文本信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = ToUser
	tmp["sendToType"] = SendToType
	tmp["sendMsgType"] = "TextMsg"
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = Qqid
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func SendPrivate(ToUser int64, Content string, PriQQid int64) {
	//发送文本信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = PriQQid
	tmp["sendToType"] = 3
	tmp["sendMsgType"] = "TextMsg"
	tmp["content"] = Content
	tmp["groupid"] = ToUser
	tmp["atUser"] = 0
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
func sendA(ToUser int64, SendToType int, Content string, SendMsgType string) {
	//发送其他信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = ToUser
	tmp["sendToType"] = SendToType
	tmp["sendMsgType"] = SendMsgType
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = 0

	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
func sendVoice(ToUser int64, SendToType int, Content string) {
	//发送语音信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = ToUser
	tmp["sendToType"] = SendToType
	tmp["sendMsgType"] = "VoiceMsg"
	tmp["content"] = ""
	tmp["voiceUrl"] = "https://dds.dui.ai/runtime/v1/synthesize?voiceId=qianranfa&speed=0.7&volume=100&audioType=wav&text=" + url.PathEscape(Content)
	tmp["groupid"] = 0
	tmp["atUser"] = 0
	tmp["voiceBase64Buf"] = ""

	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
func zan(qq1 int, err error) {
	//名片点赞
	tmp := make(map[string]interface{})
	tmp["UserID"] = qq1

	tmp1, _ := json.Marshal(tmp)
	fmt.Println(string(tmp1))
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=QQZan&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
func getinfo(qq1 int) string {
	tmp := make(map[string]interface{})
	tmp["UserID"] = qq1

	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+domain+"/v1/LuaApiCaller?funcname=GetUserInfo&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return "err"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
func getGroupUserList(qqGroupId int64, lastUin int) (*MemberData, error) {
	tmp := make(map[string]interface{})
	tmp["GroupUin"] = qqGroupId
	tmp["LastUin"] = lastUin

	tmp1, _ := json.Marshal(tmp)
	url2 := fmt.Sprintf("http://%s/v1/LuaApiCaller?funcname=GetGroupUserList&timeout=10&qq=%s", domain, qq)
	resp, err := http.Post(url2, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	v := new(MemberData)

	fmt.Println(string(body))
	err = json.Unmarshal(body, &v)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	if v == nil {
		return nil, errors.New("err")
	}

	if v.LastUin != 0 {
		// time.Sleep(1 * time.Second)
		v1, err := getGroupUserList(qqGroupId, v.LastUin)
		if err != nil {
			return nil, err
		}
		for _, member := range v1.MemberList {
			v.MemberList = append(v.MemberList, member)
		}
		v.Count += v1.Count
	}
	return v, nil
}

func getGroupList(qqGroupId int64, nextToken string) (*GroupData, error) {
	tmp := make(map[string]interface{})
	tmp["NextToken"] = nextToken
	tmp1, _ := json.Marshal(tmp)
	url2 := fmt.Sprintf("http://%s/v1/LuaApiCaller?funcname=GetGroupList&timeout=10&qq=%s", domain, qq)
	resp, err := http.Post(url2, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	v := new(GroupData)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("err")
	}
	if v.NextToken != "" {
		v1, err := getGroupList(qqGroupId, v.NextToken)
		if err != nil {
			return nil, err
		}
		for _, group := range v1.GroupList {
			v.GroupList = append(v.GroupList, group)
		}
		v.Count += v1.Count
	}
	return v, nil
}

func getGroupData(qqGroupId int64) *Group {
	qqGroupData, err := getGroupList(qqGroupId, "")
	if err != nil {
		return nil
	}
	for _, data := range qqGroupData.GroupList {
		if data.GroupId == qqGroupId {
			return &data
		}
	}
	return nil
}

func getGroupUserAutoRemark(qqGroupId int64, qqid int64) string {

	// "owner"、"admin"、"member"
	qqGroupMemberData, err := getGroupUserList(qqGroupId, 0)
	if err != nil {
		return ""
	}
	for _, data := range qqGroupMemberData.MemberList {
		if data.Qqid == qqid {
			return data.AutoRemark
		}
	}
	return ""
}
func getGroupUserData(qqGroupId int64, qqid int64) *Member {

	qqGroupMemberData, err := getGroupUserList(qqGroupId, 0)
	if err != nil {
		return nil
	}
	for _, data := range qqGroupMemberData.MemberList {
		if data.Qqid == qqid {
			return &data
		}
	}
	return nil
}
