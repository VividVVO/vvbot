package plugins

// 消息数据接收
type MeassageData struct {
	// 收到的消息
	Content string `json:"Content"`
	// 消息来自的群号
	FromGroupID int64 `json:"FromGroupId"`
	// 消息来自的群名
	FromGroupName string `json:"FromGroupName"`
	// 发送消息用户的别名
	FromNickName string `json:"FromNickName"`
	// 发送消息用户的QQ号
	FromUserID int64 `json:"FromUserId"`

	// 发送消息的来源ID(qq号或者群号)
	FromSourceID int64 `json:"FromSourceID"`

	MsgRandom int `json:"MsgRandom"`
	MsgSeq    int `json:"MsgSeq"`
	// 消息时间
	MsgTime int64 `json:"MsgTime"`
	// 消息类型
	MsgType string `json:"MsgType"`
	// 红包信息
	RedBaginfo interface{} `json:"RedBaginfo"`
	// 接收消息的人QQ号（机器人）
	CurrentQQ int64 `json:"CurrentQQ"`
	// @的QQ号
	AtQQList []int64
	// 消息来自的QQ号
	FromUin int64 `json:"FromUin"`
	// 消息接收的QQ号
	ToUin int64 `json:"ToUin"`

	// 发送消息的类型 1 私聊消息、2 群聊消息
	SendToType int `json:"SendToType"`
}

// 所有群成员信息
type MemberData struct {
	Count      int      `json:"Count"`
	LastUin    int      `json:"LastUin"`
	MemberList []Member `json:"MemberList"`
}

// 群成员信息
type Member struct {
	Qqid          int64  `json:"MemberUin"`
	Sex           string `json:"sex"` // "male"、"female"、"unknown"
	Age           int    `json:"age"`
	JoinTime      int64  `json:"JoinTime"`
	AutoRemark    string `json:"AutoRemark"`
	GroupCard     string `json:"GroupCard"`
	Role          string `json:"Role"` // "owner"、"admin"、"member"
	LastSpeakTime int64  `json:"LastSpeakTime"`
	NickName      string `json:"NickName"`
}

type ExecutePlugin interface {
	// 手动校验消息数据，需要关闭自动校验
	CheckOrder(string) bool
	// 开启关闭自动校验
	IsNotCheckOrder() bool
	// 获取命令列表
	GetOrders() []string

	/*	mess.Content 消息内容 string
		mess.FromGroupID 来源QQ群 int
		mess.FromUserID 来源QQ int64
		mess.iotqqType 消息类型 string  */
	Run(MeassageData, string, int64)
}
