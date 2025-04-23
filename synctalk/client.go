package synctalk

import (
	"8hfinal/utils"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// 加入聊天室
func JoinChatRoom(ip string) {
	conn, err := net.Dial("tcp", ip+":8080")
	if err != nil {
		log.Println("连接失败:", err)
		return
	}
	log.Println("成功链接到服务器")

	go func(conn2 net.Conn) {
		chatroom := NewChat()
		chatroom.Join(conn2)
	}(conn)

}

// 进入和服务器的私人聊天
func StartPrivateClient(ip string, port int) {
	//建立和服务器的链接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println("连接失败:", err)
	}
	log.Println("成功链接....")
	defer conn.Close()

	//读取服务端发送信息
	go func() {
		for {
			msg, err := utils.ReadMessages(conn)
			if err != nil {
				log.Println("连接断开")
				return
			}
			log.Println("收到服务器发来的消息：", msg)
		}
	}()

	//给服务器发送消息
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		text := scanner.Text()
		conn.Write([]byte(text + "\n"))
		log.Println("客户端已经给服务器发出消息", text)
	}
}

// 选择客户端，需要返回一个newserver，要输入进入的ip和port
func ClientChose() {
	var ip string
	var port int
	var chose int
START:
	log.Println("请选择您想要进行聊天的方式： 1---私人聊天 2---聊天室")
	fmt.Scanln(&chose)
	for {
		switch chose {
		case 1:

			log.Println("请输入您想要访问的服务器地址：")
			fmt.Scanln(&ip)
			StartPrivateClient(ip, port)

		case 2:
			log.Println("请输入您想要访问的服务器地址：")
			fmt.Scanln(&ip)
			log.Println("请输入您要访问的端口")
			fmt.Scanln(&port)
			JoinChatRoom(ip)
		default:
			goto START
		}

	}
}
