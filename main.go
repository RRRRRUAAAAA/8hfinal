package main

import (
	"8hfinal/synctalk"
	"fmt"
	"log"
)

//TIp <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//对最简单的服务器的测试
	//server := NewServer("127.0.0.1", 8080)
	//go server.Start()
	//testClient()

	fmt.Printf("选择模式 ：1---服务端 ， 2---客户端  ")
	var mode int
	var ip string
	var port int
	fmt.Scanln(&mode)
	switch mode {

	case 1: //服务端启动
		{
			server := synctalk.NewServer("0.0.0.0", 8080)
			server.StartPrivateServer()
		}
	case 2: //私人链接

		{
		Start1:
			for {
				var chose int
				log.Printf("请选择你的连接方式：1---进入聊天室  2---私人链接 \n")
				fmt.Scanln(&chose)
				if !(chose != 1 || chose != 2) {
					log.Println("输入有误，请重新输入：")
					goto Start1
				}
				chatroom := synctalk.NewChat()

				log.Println("请输入服务器的ip地址：\n")
				fmt.Scanln(&ip)
				log.Println("请输入你要链接的端口号\n")
				fmt.Scanln(&port)

				server := synctalk.NewServer(ip, port)
				server.DifferentServer(chatroom, chose)

			}
		}
	}
}
