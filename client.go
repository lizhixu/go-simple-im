package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	flag       int
}

var (
	serverIp   string
	serverPort int
)

func init() {
	flag.StringVar(&serverIp, "i", "127.0.0.1", "IP地址")
	flag.IntVar(&serverPort, "p", 888, "端口")
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{ServerIp: serverIp, ServerPort: serverPort, flag: 999}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("服务器连接失败，请重试", err)
		return nil
	}
	client.Conn = conn
	return client
}

func sendMenuTips() {
	fmt.Println("请按照提示选择\n 1.公聊模式\n 2.私聊模式\n 3.修改用户名\n 0.退出")
}

//菜单选项
func (this *Client) menu() bool {
	//输入字符串退出的情况
	var flagTmp string
	fmt.Scanln(&flagTmp)
	flag, err := strconv.Atoi(flagTmp)
	if err != nil {
		sendMenuTips()
		return false
	}

	if flag >= 0 && flag < 4 {
		this.flag = flag
		return true
	} else {
		sendMenuTips()
		return false
	}
}

//改名
func (client *Client) renamed() bool {
	fmt.Scanln(&client.Name)
	sendMsg := "改名|" + client.Name + "\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("改名失败", err)
		return false
	}
	return true
}

// PublicChat PublicChar 公聊模式
func (client *Client) PublicChat() {
	var charMsg string
	fmt.Scanln(&charMsg)
	for charMsg != "exit" {
		if charMsg != "" {
			sendMsg := charMsg + "\n"
			_, err := client.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("消息发送失败", err)
				break
			}
			charMsg = ""
			fmt.Scanln(&charMsg)
		}
	}
	sendMenuTips()
}

//查询当前在线用户
func (client *Client) selectUser() {
	sendMsg := "who\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("当前用户查询失败", err)
	}
}

// PriviteChat 私聊消息
func (client Client) PriviteChat() {
	fmt.Println("====请输入对方用户名，exit退出")
	var otherName string
	fmt.Scanln(&otherName)
	for otherName != "exit" {
		fmt.Println("====私聊会话已接通，exit退出")
		var charMsg string
		fmt.Scanln(&charMsg)
		for charMsg != "exit" {
			if charMsg != "" {
				sendMsg := "@" + otherName + " " + charMsg + "\n"
				_, err := client.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("消息发送失败", err)
					break
				}
				charMsg = ""
				fmt.Scanln(&charMsg)
			}
		}
	}
}

// Run 运行菜单
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("====进入公聊模式，exit退出\n")
			client.PublicChat()
			break
		case 2:
			fmt.Println("====进入私聊模式\n")
			//展示当前用户
			client.selectUser()
			client.PriviteChat()
			break
		case 3:
			fmt.Println("====请输入用户名\n")
			client.renamed()
			fmt.Println("====改名成功，请按照菜单继续操作")
			break
		}
	}
}

// ListenResponse 处理server消息的回执
func (client *Client) ListenResponse() {
	//直接将消息发送至控制台
	io.Copy(os.Stdout, client.Conn)
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("===服务器连接失败")
		return
	} else {
		fmt.Println("===服务器连接成功")
		go client.ListenResponse()
		fmt.Println("=====请先输入一个中意昵称")
		client.renamed()
		sendMenuTips()
	}
	client.Run()
}
