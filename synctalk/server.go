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

func NewServer(ip string, port int) *Server {
	s := new(Server)
	s.Port = port
	s.Ip = ip
	return s
}
func Handle(conn net.Conn) {
	defer conn.Close()
	log.Println("客户端已连接 :", conn.RemoteAddr())
	fmt.Println("服务器已经成功链接")

	//发送信息
	go func() {
		for {
			msg, err := utils.ReadMessages(conn)
			if err != nil {
				log.Println("客户端断开：", err)
				return
			}
			log.Println("客户端说：", msg)
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
func (s *Server) StartServer() {
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
		go Handle(conn)
	}

}
