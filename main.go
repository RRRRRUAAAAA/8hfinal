package main

import (
	synctalk2 "8hfinal/failfunction/synctalk"
	"fmt"
	"log"
)

//TIp <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	var chose int
	log.Println("选择您的主机类型 ： 1---服务器  2---客户端")
	fmt.Scanln(&chose)
	switch chose {
	case 1:
		server := synctalk2.NewServer("0.0.0.0", 8080)
		server.DifferentServer()
	case 2:
		synctalk2.ClientChose()

	}
}
