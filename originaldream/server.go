package originaldream

import (
	"8hfinal/utils"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 开辟协程帮助处理连接
func (this *Server) Handle(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("客户端的地址是：%v\n", conn.RemoteAddr())
	//发送消息给客户端
	_, _ = conn.Write([]byte("欢迎来到服务器\n"))
	//接收客户端的消息
	//buf := make([]byte, 1024)
	//n, err := conn.Read(buf)
	//if err != nil {
	//	log.Println("读取失败", err)
	//	return
	//}
	//recv := string(buf[:n])
	recv, _ := utils.ReadMessages(conn)
	fmt.Printf("我是服务器，我收到了客户端发来的：%s\n", recv)

	//回复客户端
	reply := "我是客户端，已经收到你发的" + recv
	conn.Write([]byte(reply))
	//链接成功
	fmt.Println("客户端链接成功：", conn.RemoteAddr())

}
func NewServer(Ip string, Port int) *Server {
	s := new(Server)
	s.Port = Port
	s.Ip = Ip
	return s
}

func (this *Server) Start() {
	//建立监听器
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		log.Println("监听器设置出错")
	}
	log.Println("服务器启动，监听🀄️")
	defer listener.Close()
	//建立连接
	for {

		conn, err := listener.Accept()

		fmt.Println("服务已经启动，等待客户连接。。。")
		if err != nil {
			log.Printf("出现错误 ：%v", err)
			continue

		}
		go this.Handle(conn)
	}

}
