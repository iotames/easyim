package main

import (
	"flag"

	"github.com/iotames/easyim/config"
	"github.com/iotames/easyim/database"
	"github.com/iotames/easyim/server"
)

var (
	appInit bool
)

func main() {
	sconf := config.GetServer()
	flag.IntVar(&sconf.Port, "port", sconf.Port, "监听端口")

	flag.Parse()

	if appInit {
		database.SyncTables()
	}
	server := server.NewServer(sconf)
	server.Start()
}

func init() {
	config.LoadEnv()

	flag.BoolVar(&appInit, "init", false, "首次运行时添加，用于初始化")
	// time.LoadLocation("Asia/Shanghai")
}
