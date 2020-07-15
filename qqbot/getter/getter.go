package getter

// 消息数据接收
type MeassageData struct {
	// 收到的消息
	Content string `json:"Content"`
	// 消息来自的群号
	FromGroupID int `json:"FromGroupId"`
	// 消息来自的群名
	FromGroupName string `json:"FromGroupName"`
	// 发送消息用户的别名
	FromNickName string `json:"FromNickName"`
	// 发送消息用户的QQ号
	FromUserID int64 `json:"FromUserId"`

	MsgRandom int `json:"MsgRandom"`
	MsgSeq    int `json:"MsgSeq"`
	// 消息时间
	MsgTime int `json:"MsgTime"`
	// 消息类型
	MsgType string `json:"MsgType"`
	// 红包信息
	RedBaginfo interface{} `json:"RedBaginfo"`
	// 接收消息的人QQ号（机器人）
	CurrentQQ int64 `json:"CurrentQQ"`
	// @的QQ号
	AtQQList []int64
}

// 所有群成员信息
type MemberData struct {
	Count      int      `json:"Count"`
	LastUin    int      `json:"LastUin"`
	MemberList []Member `json:"MemberList"`
}

// 群成员信息
type Member struct {
	JoinTime      int64  `json:"JoinTime"`
	AutoRemark    string `json:"AutoRemark"`
	GroupCard     string `json:"GroupCard"`
	LastSpeakTime int64  `json:"LastSpeakTime"`
	Qqid          int    `json:"MemberUin"`
	NickName      string `json:"NickName"`
}

type ExecuteOrder interface {
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
	Run(MeassageData, string, int)
}
