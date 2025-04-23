package synctalk

import (
	"8hfinal/utils"
	"bufio"
	"log"
	"net"
	"os"
)

func StartClient(ip string) {
	//建立和服务器的链接
	conn, err := net.Dial("tcp", ip+":8080")
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
		text := scanner.Text()
		conn.Write([]byte(text + "\n"))
		log.Println("客户端已经给服务器发出消息", text)
	}
}
