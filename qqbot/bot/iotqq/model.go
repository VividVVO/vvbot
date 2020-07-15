package iotqq

import (
	"bytes"
	"encoding/json"
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
type Member struct {
	JoinTime      int64  `json:"JoinTime"`
	AutoRemark    string `json:"AutoRemark"`
	GroupCard     string `json:"GroupCard"`
	LastSpeakTime int64  `json:"LastSpeakTime"`
	Qqid          int    `json:"MemberUin"`
	NickName      string `json:"NickName"`
}

func getQQ() string {
	return qq
}

func sendPic(ToUser int, SendToType int, Content string, PicBase64Buf string, PicUrl string) {
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
func send(ToUser int, SendToType int, Content string, Qqid int) {
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
func sendA(ToUser int, SendToType int, Content string, SendMsgType string) {
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
func sendVoice(ToUser int, SendToType int, Content string) {
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
func getGroupUserList(qqGroupId int, lastUin int) (*MemberData, error) {
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
		return nil, err
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
func getGourpUserAutoRemark(qqGroupId int, qqid int) string {
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
func getGourpUserData(qqGroupId int, qqid int) *Member {
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
