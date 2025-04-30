package channelchatroom

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Status struct {
}
type User struct {
	Addr   string //使用客户端的地址
	status Status
	Name   string
	C      chan string
}

// 在server建立连接的时候调用直接创建,同时创建一个服务端
func NewUser() *User {
	return &User{
		C: make(chan string),
	}
}

func (u *User) AskUserdetail(conn net.Conn) {
	fmt.Println("请输入你的用户名，不输入的话，默认为匿名")
	conn.Write([]byte("请输入您的用户名"))
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			log.Println("用户自己断开连接")
		} else {
			log.Println("连接出现错误（用户输入用户名时）")
		}

	}
	name := string(buf[:n])
	name = strings.TrimSpace(name)
	u.Name = name
	fmt.Printf("您的用户名为：%s\n", u.Name)
	addr := conn.RemoteAddr()
	u.Addr = addr.String()
}
func (u *User) Listener() {
	for msg := range u.C {
		fmt.Println("[服务器消息]", msg)
	}
}

func (u *User) writter(conn net.Conn) {
	var msg string
	fmt.Println("请输入您想发送的内容：")
	fmt.Println("1:直接发送表示聊天室广播消息")
	fmt.Println("2:@xxx表示私聊消息")
	fmt.Println("3:/指令表示发送系统指令操作，如果想查看全部指令，请输入/all")
	for {
		fmt.Scanln(&msg)
		_, err := conn.Write([]byte(msg))
		if err != nil {
			if err == io.EOF {
				log.Println("用户正常退出")
			} else {
				log.Println("连接崩溃")
			}
			break
		}
	}
}
func (u *User) Start(conn net.Conn) {
	// 启动后台监听协程，从 u.C 中读取并打印消息
	go u.Listener()

	// 主线程运行输入写入逻辑
	u.writter(conn)
}
