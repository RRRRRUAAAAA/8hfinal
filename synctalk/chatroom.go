package synctalk

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type ChatRoom struct {
	clients map[net.Conn]string //key = conn 标识唯一用户 value = 用户名
	mu      sync.Mutex          //锁保证并发安全
}

func NewChat() *ChatRoom {
	return &ChatRoom{
		clients: make(map[net.Conn]string)}
}

func (chat *ChatRoom) Join(conn net.Conn) {

	name := chat.AskName(conn) //询问姓名
	chat.mu.Lock()
	chat.clients[conn] = name //为新来的用户赋值姓名
	chat.mu.Unlock()

	chat.BroadCast(conn, fmt.Sprintf("%s 加入聊天室", name)) //广播告知该用户加入聊天室

	go chat.HandleMessage(conn, name)
}

// 用来得到用户输入的名字
func (chat *ChatRoom) AskName(conn net.Conn) string {
	conn.Write([]byte("请输入你的用户名"))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		text := scanner.Text()
		return text

	}
	return "匿名"
}

// 处理消息
func (chat *ChatRoom) HandleMessage(conn net.Conn, name string) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("chatroom中conn流关闭出现错误")
		}
	}(conn)
Start1:
	conn.Write([]byte("请输入您想输入的内容"))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		msg := scanner.Text()
		if strings.HasPrefix(msg, "@") {
			//私聊逻辑
			split := strings.SplitN(msg, " ", 2)
			if len(split) < 2 {
				conn.Write([]byte("输入有误，正确的格式是: @用户名 内容\n"))
				goto Start1
			}
			targetName := strings.TrimPrefix(split[0], "@")
			chat.PrivateMessage(conn, fmt.Sprintf("[私聊] 【%s】对你说：%s", name, split[1]), targetName)
			log.Println("私聊信息已发送...")
		} else {
			chat.BroadCast(conn, fmt.Sprintf("[%s]: %s\n", name, msg))
			log.Println("广播信息已经投放...")
		}

	}
	//sanner.Scan()如果断开，意味着用户已经断开
	chat.mu.Lock()
	delete(chat.clients, conn)
	chat.mu.Unlock()
	chat.BroadCast(conn, fmt.Sprintf("[%s] 离开了聊天室", name))

}

// 全体广播
func (chat *ChatRoom) BroadCast(sender net.Conn, msg string) {
	chat.mu.Lock()
	defer chat.mu.Unlock()
	for conn := range chat.clients {
		if conn != sender {
			conn.Write([]byte(msg))
			log.Printf("已发送消息给%s", chat.clients[conn])
		}
	}
}

// 私聊
func (chat *ChatRoom) PrivateMessage(sender net.Conn, msg string, toName string) {
	chat.mu.Lock()
	defer chat.mu.Unlock()
	found := false
	for conn, name := range chat.clients {
		if name == toName {
			conn.Write([]byte(msg))
			found = true
			log.Println("发送成功")
			break
		}
	}
	if found == false {
		sender.Write([]byte("没有用户或者该用户未上线"))
		log.Printf("发送私聊消息失败")

	}
}
