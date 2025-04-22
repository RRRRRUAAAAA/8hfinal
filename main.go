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

	fmt.Printf("选择模式 ：1---服务端 ，2---客户端")
	var mode int
	var ip string
	fmt.Scanln(&mode)
	for {
		if mode == 1 {
			server := synctalk.NewServer("0.0.0.0", 8080)
			server.StartServer()
		} else if mode == 2 {
			fmt.Scanln(&ip)
			synctalk.StartClient(ip)
		} else {
			log.Println("您输入的值是非法的，麻烦输入1或2")
		}
	}
}
