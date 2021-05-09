package main

import (
	"net"
	"strings"
)

type User struct {
	Name, Addr string
	C          chan string
	Conn       net.Conn
	Server     *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		Conn:   conn,
		Server: server,
	}
	go user.ListenMessage()
	return user
}

// UserOnline 用户上线
func (this *User) UserOnline() {
	this.Server.MapLock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.MapLock.Unlock()
	//推送上线消息
	this.Server.BroadCast(this, "已上线\n")
}

// UserOffline 用户下线
func (this *User) UserOffline() {
	this.Server.MapLock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.MapLock.Unlock()
	//推送上线消息
	this.Server.BroadCast(this, "已下线\n")
}

// DoMessage 发送消息
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.Server.MapLock.Lock()
		userMsg := "当前在线用户："
		for _, user := range this.Server.OnlineMap {
			userMsg += user.Name + ","
		}
		this.SendMsg(strings.Trim(userMsg, ",") + "\n")
		this.Server.MapLock.Unlock()
	} else if len(msg) > 1 && msg[:1] == "@" {
		//私信
		msgRes := strings.Split(msg, " ")
		toUse, ok := this.Server.OnlineMap[msgRes[0][1:]]
		if !ok {
			this.Server.BroadCast(this, msg+"\n")
		} else {
			content := msgRes[1]
			if content == "" {
				this.Server.BroadCast(this, msg+"\n")
			} else {
				toUse.SendMsg(this.Name + "的私信：" + content + "\n")
			}
		}
	} else if len(msg) > 7 && msg[:7] == "改名|" {
		newName := msg[7:]
		_, ok := this.Server.OnlineMap[newName]
		if ok {
			this.SendMsg("用户名已存在\n")
		} else {
			this.Server.MapLock.Lock()
			delete(this.Server.OnlineMap, this.Name)
			this.Name = newName
			this.Server.OnlineMap[newName] = this
			this.SendMsg("您的用户名已更新为：" + newName + "\n")
			this.Server.MapLock.Unlock()
		}
	} else {
		this.Server.BroadCast(this, msg+"\n")
	}
}

// SendMsg 给指定用户发送消息
func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg))
}
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		_, err := this.Conn.Write([]byte(msg))
		if err != nil {
			return
		}
	}
}
