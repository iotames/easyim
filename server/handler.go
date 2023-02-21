package server

import (
	"fmt"
	"net"
	"time"

	"github.com/iotames/easyim/server/user"
)

func Handler(s *Server, conn net.Conn) {
	//...当前链接的业务
	fmt.Println("链接建立成功")
	u := user.NewUser(conn, s)
	u.ConnectStart()

	//接受客户端发送的消息
	go func() {
		for {
			err := u.MsgHandler(conn)
			if err != nil {
				return
			}
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-u.GetActiveChannel():
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器

		case <-time.After(time.Second * time.Duration(s.DropAfter)):
			//已经超时
			//将当前的User强制的关闭
			u.Close()
			//退出当前Handler
			return //runtime.Goexit()
		}
	}
}
