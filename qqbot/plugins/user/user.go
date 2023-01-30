package pcr

import (
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"github.com/vivid-vvo/vvbot/qqbot/plugins/user/plugin"
)

// 添加需要的命令 (群聊)
var ExecuteGroupPlugins = []plugins.ExecutePlugin{}

// 添加需要的命令 (私聊)
var ExecuteFriendPlugins = []plugins.ExecutePlugin{
	plugin.Denglu{},
	plugin.Suijimima{},
}
