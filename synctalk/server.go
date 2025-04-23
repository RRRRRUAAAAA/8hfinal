package synctalk

import (
	"8hfinal/utils"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var privateServerStarted bool
var mu sync.Mutex
var chatroomStarted bool

type Server struct {
	Ip   string
	Port int
}

func (s *Server) DifferentServer() {
	var chose int
START:
	log.Println("请选择服务器的程序： 1---私人聊天 2---聊天室")
	fmt.Scanln(&chose)
	for {

		switch chose {
		case 1:
			//私人聊天：
			mu.Lock()
			go func() {
				if privateServerStarted {
					log.Println("私人服务器已经启动")
					mu.Unlock()
				}
				privateServerStarted = true
				s.StartPrivateServer()
			}()

		case 2:
			//聊天室
			go func() {
				mu.Lock()
				if chatroomStarted {
					log.Println("私人服务器已经启动")
					mu.Unlock()
				}
				chatroomStarted = true
				chat := NewChat()
				s.StartChatRoom(chat)

			}()
		default:
			goto START
		}

	}

}

func NewServer(ip string, port int) *Server {
	s := new(Server)
	s.Port = port
	s.Ip = ip
	return s
}
func HandlePrivate(conn net.Conn) {
	defer conn.Close()
	log.Println("私人客户端已连接 :", conn.RemoteAddr())
	fmt.Println("私人服务器已经成功链接")

	//发送信息
	go func() {
		for {
			msg, err := utils.ReadMessages(conn)
			if err != nil {
				log.Println("私人客户端断开：", err)
				return
			}
			log.Println("私人客户端说：", msg)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		text := scanner.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Println("发送失败，错误：", err)
			return
		}
	}
}

// 开启聊天室
func (s *Server) StartChatRoom(chatroom *ChatRoom) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("聊天室监听失败:", err)
		return
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {

		}
	}(ln)
	log.Println("聊天室服务已经启动")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("聊天室连接错误:", err)
			continue
		}
		log.Println("💬 新用户加入聊天室：", conn.RemoteAddr())
		go chatroom.Join(conn)
	}

}

// 开启私人服务器
func (s *Server) StartPrivateServer() {
	//建立监听器

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			log.Println("监听器出现错误：", err)
			return
		}
	}(ln)
	if err != nil {
		log.Println("监听器出现错误：", err)
		return
	} else {
		log.Println("监听器正在监听...")
	}
	//阻塞等待客户端链接
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("链接出现问题：", err)
			return
		} else {
			log.Println("链接成功")
		}
		go HandlePrivate(conn)
	}

}
