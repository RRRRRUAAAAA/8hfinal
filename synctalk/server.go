package synctalk

import (
	"8hfinal/utils"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	Ip   string
	Port int
}

func (s *Server) DifferentServer(chat *ChatRoom, chose int) {

	switch chose {
	case 1:
		//私人聊天：
		go s.StartPrivateServer()
	case 2:
		//聊天室
		go s.StartChatRoom(chat)
	}

	//堵塞主线程
	select {}
}
func NewServerEmpty() *Server {
	s := new(Server)
	return s
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
	ln, err := net.Listen("tcp", fmt.Sprintf("%s,%d", s.Ip, s.Port))
	if err != nil {
		log.Println("聊天室监听失败:", err)
		return
	}
	defer ln.Close()
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
	defer ln.Close()
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
