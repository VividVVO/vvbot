package pcr

import (
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/getter/pcr/order"
)

// 添加需要的命令
var ExecuteOrders = []getter.ExecuteOrder{
	order.Suijisetu{},
	order.Shenqchudao{},
	order.Jiesuo{},
	order.Quxiaochudao{},
	order.Chaungjiangonghui{},
	order.Jiebanggonghui{},
	order.Bangdinggonghui{},
	order.Jiarugonghui{},
	order.Tuichugonghui{},
	order.Gonghuiliebiao{},
	order.Kaiqigonghuizhan{},
	order.Baodao{},
	order.Chedao{},
	order.Jiaruquanbuchengyuan{},
	order.Zhuangtai{},
	order.Sl{},
	order.Guashu{},
	order.Suoding{},
	order.Quxiaochedao{},
	order.Xiuzhengboss{},
	order.Chexiao{},
	order.Budao{},
	order.Chadao{},
}
