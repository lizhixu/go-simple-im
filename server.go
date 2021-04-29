package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Point     int
	OnlineMap map[string]*User
	MapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, point int) *Server {
	server := &Server{
		Ip:        ip,
		Point:     point,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Start() {
	//监听socket
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Point))
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
	user := NewUser(conn)
	this.MapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.MapLock.Unlock()
	//推送上线消息
	this.BroadCast(user, "已上线")
}
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
