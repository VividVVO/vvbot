package order

import (
	"encoding/base64"
	"fmt"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
)

type Suijisetu struct {
}

func (l Suijisetu) CheckOrder(cm string) bool {
	return true
}

func (l Suijisetu) IsNotCheckOrder() bool {
	return false
}

func (l Suijisetu) GetOrders() []string {
	return []string{
		"^(随|谁)机(色|蛇)图$",
		"^申请色图$",
		"^我要看色图$",
		"^我要色图$",
		"^更多色图$",
		"^看看色图$",
		"^来点色图$",

		"^发(.+)张色图$",
		"^色图$",
	}
}

func (l Suijisetu) Run(mess getter.MeassageData, cm string, atqq int) {
	//filePth := "C:/Users/qq112/Desktop/ps_data/pcr/色图/78690381_p0.png"
	filePth := "C:/Users/qq112/Desktop/ps_data/pcr/色图/陈睿.png"
	input, err := ioutil.ReadFile(filePth)
	if err != nil {
		fmt.Println("read fail", err)
		return
	}
	// base64.StdEncoding.EncodeToString(buf)
	/*
		req, _ := http.NewRequest("GET", "https://i.pximg.net/img-original/img/2020/07/11/00/09/53/82887547_p0.png", nil)
		req.Header.Set("referer", "https://www.pixiv.net/")
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			glog.Fatal(err)
			return
		}
		var (
			img1 image.Image
		)
		defer resp.Body.Close()
		if img1, _, err = image.Decode(resp.Body); err != nil {
			glog.Error(err)
			return
		}
		/*b := img1.Bounds()
		width := b.Max.X
		height := b.Max.Y*/
	/*m1 := resize.Resize(uint(600), 0, img1, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, m1)
	if err != nil {
		glog.Error(err)
		return
	}*/
	bot.SendPic(mess.FromGroupID, 2, "", base64.StdEncoding.EncodeToString(input), "")
}
