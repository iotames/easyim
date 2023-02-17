package server

import (
	"fmt"
	"io"
	"net"
	"time"
)

func Handler(s *Server, conn net.Conn) {
	//...当前链接的业务
	fmt.Println("链接建立成功")

	user := NewUser(conn, s)
	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			//用户针对msg进行消息处理
			user.DoMessage(msg)

			//用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器

		case <-time.After(time.Second * 300):
			//已经超时
			//将当前的User强制的关闭

			user.SendMsg("你被踢了")

			//销毁用的资源
			close(user.Message)

			//关闭连接
			conn.Close()

			//退出当前Handler
			return //runtime.Goexit()
		}
	}
}
