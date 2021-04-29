package main

import "net"

type User struct {
	Name, Addr string
	C          chan string
	Conn       net.Conn
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name: addr,
		Addr: addr,
		C:    make(chan string),
		Conn: conn,
	}
	go user.ListenMessage()
	return user
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
