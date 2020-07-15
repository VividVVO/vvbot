package boss_data

import (
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"io/ioutil"
)

type BossData struct {
	Jp  [][]int     `json:"jp"`
	Cn  [][]int     `json:"cn"`
	Tw  [][]int     `json:"tw"`
	Eff [][]float32 `json:"eff"`
}

var bossData *BossData

func Init() {
	bossData = new(BossData)
	if err := load("./boss3.json", &bossData); err != nil {
		glog.Error("boss数据加载失败，请检查目录下 boss3.json 文件是否缺失")
		return
	}
}

func GetBossHpList(gameName string, cycle int) []int {
	var hpList []int
	switch gameName {
	case "JP":
		if cycle > len(bossData.Jp) {
			cycle = len(bossData.Jp)
		}
		hpList = bossData.Jp[cycle-1]
	case "TW":
		if cycle > len(bossData.Tw) {
			cycle = len(bossData.Tw)
		}
		hpList = bossData.Tw[cycle-1]
	case "KR":
		if cycle > len(bossData.Tw) {
			cycle = len(bossData.Tw)
		}
		hpList = bossData.Tw[cycle-1]
	case "CN":
		if cycle > len(bossData.Cn) {
			cycle = len(bossData.Cn)
		}
		hpList = bossData.Cn[cycle-1]
	}
	return hpList
}

func load(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
