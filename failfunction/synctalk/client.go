package synctalk

import (
	"8hfinal/failfunction/utils"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// 加入聊天室
func JoinChatRoom(ip string) {
	//第一步建立连接
	conn, err := net.Dial("tcp", ip+":8080")
	if err != nil {
		log.Println("加入聊天室出线错误：", err)
	}
	log.Println("成功加入聊天室")

	//第二步 接收聊天室发送的消息并且设置用户名
	readScanner := bufio.NewScanner(conn)
	if readScanner.Scan() {
		fmt.Println(readScanner.Text())
	}
	writeSacnner := bufio.NewScanner(os.Stdin)
	if writeSacnner.Scan() {
		name := writeSacnner.Text()
		conn.Write([]byte(name + "\n"))
	}

	//第三步 开启协程一直监听聊天室

	go func() {
		for {
			if readScanner.Scan() {
				msg := readScanner.Text()
				fmt.Printf("%s\n", msg)
			}
		}
		log.Println("服务器连接断开或出现错误")
	}()

	//第四步 始终维持用户键盘的监听
	for {
		if writeSacnner.Scan() {
			write := writeSacnner.Text()
			if write != "" {
				_, err := conn.Write([]byte(write + "\n"))
				if err != nil {
					log.Println("发送失败，断开连接", err)
					return
				}
			}

		}
	}
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

	switch chose {

	case 1:
		for {
			log.Println("请输入您想要访问的服务器地址：")
			fmt.Scanln(&ip)
			log.Println("请输入您要访问的端口")
			fmt.Scanln(&port)
			StartPrivateClient(ip, port)
		}
	case 2:
		log.Println("请输入您想要访问的服务器地址：")
		fmt.Scanln(&ip)
		JoinChatRoom(ip)
	default:
		goto START

	}
}
