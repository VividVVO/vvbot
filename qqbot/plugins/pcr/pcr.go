package pcr

import (
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"github.com/vivid-vvo/vvbot/qqbot/plugins/pcr/plugin"
)

// 添加需要的插件 (群聊)
var ExecuteGroupPlugins = []plugins.ExecutePlugin{
	plugin.Suijisetu{},
	plugin.Shenqchudao{},
	plugin.Jiesuo{},
	plugin.Quxiaochudao{},
	plugin.Chaungjiangonghui{},
	plugin.Jiebanggonghui{},
	plugin.Bangdinggonghui{},
	plugin.Jiarugonghui{},
	plugin.Tuichugonghui{},
	plugin.Gonghuiliebiao{},
	plugin.Kaiqigonghuizhan{},
	plugin.Baodao{},
	plugin.Chedao{},
	plugin.Jiaruquanbuchengyuan{},
	plugin.Zhuangtai{},
	plugin.Sl{},
	plugin.Guashu{},
	plugin.Suoding{},
	plugin.Quxiaochedao{},
	plugin.Xiuzhengboss{},
	plugin.Chexiao{},
	plugin.Budao{},
	plugin.Chasu{},
	plugin.Quanbuxiashu{},
	plugin.Chadao{},
	plugin.Mianban{},
}

// 添加需要的命令 (私聊)
var ExecuteFriendPlugins = []plugins.ExecutePlugin{}
