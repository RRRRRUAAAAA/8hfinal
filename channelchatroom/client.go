package channelchatroom

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
	Ip   string
	Port int
	C    chan string
}

// 调用newClient同时也会创建一个user
func NewClient(ip string, port int, usr *User) *Client {

	return &Client{
		Ip:   ip,
		Port: port,
		C:    usr.C,
	}
}

func (c *Client) Connect() net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Ip, c.Port))
	if err != nil {
		if err == io.EOF {
			log.Println("连接断开", err)
		} else {
			log.Println("连接崩溃（客户端与服务端）：", err)
		}
	}
	go func() {
		c.receive(conn)
	}()
	return conn
}

func (c *Client) receive(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("用户断开连接", err)
			} else {
				log.Println("连接断开（客户端读）", err)
			}
			break
		}
		c.C <- string(buf[:n])
	}
}
