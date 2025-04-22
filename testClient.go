package main

import (
	"8hfinal/utils"
	"fmt"
	"log"
	"net"
	"time"
)

func testClient() {
	//time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("连接失败:", err)
	}

	//向客户端发送消息
	//conn.Write([]byte("我连接了你，服务器\n"))
	////接收客户端消息
	//buf := make([]byte, 1024)
	//n, _ := conn.Read(buf)
	//resMsg := string(buf[:n])
	resMsg, _ := utils.ReadMessages(conn)
	fmt.Printf("收到客户端消息：%s\n", resMsg)

	//回复客户端
	wrtMsg := "我是服务器端：我收到了你发的：" + resMsg
	conn.Write([]byte(wrtMsg))
	defer conn.Close()
	fmt.Println(conn, "Hello Server!]n")
	time.Sleep(3 * time.Second)
}
