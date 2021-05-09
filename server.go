package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port     int
	OnlineMap map[string]*User
	MapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:     port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Start() {
	//监听socket
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("监听失败", err)
		return
	}
	//关闭监听
	defer ln.Close()

	//启动监听Message的goruntime
	go this.ListenMessage()
	for {
		//等待下一个连接
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("fail")
			continue
		}
		go this.Handler(conn)
	}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")
	//上线用户加入广播
	user := NewUser(conn, this)
	//推送上线消息
	user.UserOnline()
	isActivie := make(chan bool)
	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4655)
		for {
			r, err := conn.Read(buf)
			if r == 0 {
				user.UserOffline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
			msg := string(buf[:r-1])
			user.DoMessage(msg)
			//设置活跃状态
			isActivie <- true
		}
	}()
	//当前handler阻塞
	for {
		select {
		case <-isActivie:
		case <-time.After(time.Minute * 10):
			user.SendMsg("会话超时已断开连接\n")
			close(user.C)
			conn.Close()
			return//退出当前handler
		}
	}
}

// BroadCast 向所有用户推送消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		//将消息发送给所有用户
		this.MapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.MapLock.Unlock()
	}
}
