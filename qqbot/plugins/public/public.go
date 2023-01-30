package pcr

import (
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"github.com/vivid-vvo/vvbot/qqbot/plugins/public/plugin"
)

// 添加需要的插件 (群聊)
var ExecuteGroupPlugins = []plugins.ExecutePlugin{
	plugin.Version{},
	plugin.LeaveGroup{},
	plugin.Help{},
}

// 添加需要的命令 (私聊)
var ExecuteFriendPlugins = []plugins.ExecutePlugin{
	plugin.Version{},
}
