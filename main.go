package main

import (
	"flag"
	"fmt"

	"github.com/iotames/easyim/config"
	"github.com/iotames/easyim/database"
	"github.com/iotames/easyim/server"
)

var (
	daemon  bool // 以守护进程方式后台运行
	stop    bool // 停止程序
	appInit bool
	sconf   config.Server
)

func main() {
	var err error
	flag.Parse()
	if stop {
		err = stopApp(sconf.Port)
		if err != nil {
			fmt.Println("stop err:", err)
			return
		}
		return
	}
	if daemon {
		err = startDaemon()
		if err != nil {
			fmt.Println("start daemon err:", err)
			return
		}
		return
	}
	if appInit {
		database.SyncTables()
	}
	server := server.NewServer(sconf)
	err = server.Start()
	if err != nil {
		fmt.Println("EasyIM start fail:", err)
	}
}

func init() {
	config.LoadEnv()
	flag.BoolVar(&daemon, "d", false, "以守护进程的方式在后台运行。不支持windows系统")
	flag.BoolVar(&stop, "stop", false, "停止运行中的程序")
	flag.BoolVar(&appInit, "init", false, "首次运行时添加，用于初始化")
	sconf = config.GetServer()
	flag.IntVar(&sconf.Port, "port", sconf.Port, "监听的端口号")
	// time.LoadLocation("Asia/Shanghai")
}
