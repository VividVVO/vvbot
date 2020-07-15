package main

import (
	"github.com/gogf/gf/frame/g"
	_ "github.com/vivid-vvo/vvbot/boot"
	_ "github.com/vivid-vvo/vvbot/router"
)

func main() {
	g.Server().Run()
}
