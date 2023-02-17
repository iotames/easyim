package main

import (
	"flag"
	"fmt"

	"github.com/iotames/easyim/config"
	"github.com/iotames/easyim/database"
)

var (
	serverPort int
	appInit    bool
)

func main() {
	flag.Parse()

	if appInit {
		database.SyncTables()
	}
	listenIP := "0.0.0.0"
	server := NewServer(listenIP, serverPort)
	fmt.Printf("Start EasyIM In: %s:%d\n", listenIP, serverPort)
	server.Start()
}

func init() {
	config.LoadEnv()
	sconf := config.GetServer()
	flag.IntVar(&serverPort, "port", sconf.Port, "监听端口")
	flag.BoolVar(&appInit, "init", false, "首次运行时添加，用于初始化")
	// time.LoadLocation("Asia/Shanghai")
}
