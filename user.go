package main

import "net"

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
	this.Server.BroadCast(this, msg)
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
