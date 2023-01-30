package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	_ "github.com/vivid-vvo/vvbot/boot"
	_ "github.com/vivid-vvo/vvbot/router"
)

// VERSION 版本号
var VERSION = "v1.0.0 alpha"

func main() {
	gcache.Set("VERSION", VERSION, 0)
	s := g.Server()
	s.EnableAdmin()
	s.Run()

}
