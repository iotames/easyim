package server

import (
	"fmt"
	"net"
	"time"

	"github.com/iotames/easyim/server/handler"
	"github.com/iotames/easyim/server/user"
)

// Handler 当前链接的业务
func Handler(s *Server, conn net.Conn) {
	u := user.NewUser(conn)
	u.SetOnConnectStart(func(u user.User) {
		fmt.Println("TCP连接建立成功:", conn.RemoteAddr().String())
	})
	u.SetOnConnectLost(func(u user.User) {
		addr := u.GetConn().RemoteAddr().String()
		s.UserOffline(addr)
		fmt.Println("TCP连接断开")
	})
	u.ConnectStart()

	//接受客户端发送的消息
	go func() {
		for {
			err := handler.MainHandler(s, u)
			if err != nil {
				if err.Error() == user.ERR_CONNECT_LOST {
					return
				}
				if !u.IsClosed {
					u.Close()
					fmt.Println("--conn-Closed--After--MainHandler--error:", err)
				}
				return
			}
			if !u.IsClosed {
				//用户的任意消息，代表当前用户是一个活跃的
				u.KeepActive()
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
			if !u.IsClosed {
				u.Close()
			}
			//退出当前Handler
			return //runtime.Goexit()
		}
	}
}
