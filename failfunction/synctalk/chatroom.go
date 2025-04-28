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

func (chat *ChatRoom) setName(conn net.Conn, name string) {
	chat.mu.Lock()
	chat.clients[conn] = name
	chat.mu.Unlock()
}

func (chat *ChatRoom) GetName(conn net.Conn) string {
	chat.mu.Lock()
	name := chat.clients[conn]
	chat.mu.Unlock()
	return name
}
func NewChat() *ChatRoom {
	return &ChatRoom{
		clients: make(map[net.Conn]string),
	}
}

func (chat *ChatRoom) Join(conn net.Conn) {

	name := chat.AskName(conn) //询问姓名
	chat.setName(conn, name)
	chat.BroadCast(conn, fmt.Sprintf("%s 加入聊天室\n", name)) //广播告知该用户加入聊天室
	go chat.HandleMessage(conn)

	defer log.Printf("%s成功加入聊天室", name)
}

// 用来得到用户输入的名字
func (chat *ChatRoom) AskName(conn net.Conn) string {
	conn.Write([]byte("请输入你的用户名\n"))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		text := scanner.Text()
		conn.Write([]byte(fmt.Sprintf("命名成功，您的用户名是：%s 退出指令：/exit 查询在线人数指令：/who\n 更改姓名： /rename\n", text+"\n")))
		return text

	}
	conn.Write([]byte(fmt.Sprintf("命名成功，您是匿名用户 退出指令：/exit 查询在线人数指令：/who\n 更改姓名： /rename\n" + "\n")))
	return "匿名"
}
func (chat *ChatRoom) Rename(conn net.Conn, name string) {

	chat.mu.Lock()
	oldname := chat.clients[conn]
	if chat.clients[conn] == name {
		text := "请输入和原来不同的名字"
		log.Printf("用户%s修改名字失败，原因是：%s", chat.clients[conn], text)
		conn.Write([]byte("姓名修改失败：" + text + "\n"))
		chat.mu.Unlock()
		return
	}

	for connrange := range chat.clients {
		if conn != connrange && chat.clients[connrange] == name {
			text := "该用户名已存在。"
			log.Printf("用户%s修改名字失败，原因是：%s", chat.clients[conn], text)
			conn.Write([]byte("姓名修改失败：" + text + "\n"))
			chat.mu.Unlock()
			return
		}
	}
	chat.clients[conn] = name
	conn.Write([]byte(fmt.Sprintf("用户名更改成功，您的用户名为：%s 您可以继续聊天了\n", name+"\n")))
	log.Printf("用户%s已经改名为：%s", oldname, name)
	chat.BroadCast(conn, fmt.Sprintf("用户%s已经改名为：%s ", oldname, name+"\n"))
	chat.mu.Unlock()
	return
}

// 处理消息
func (chat *ChatRoom) HandleMessage(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("chatroom中conn流关闭出现错误")
		}
	}(conn)

	conn.Write([]byte("您已经进入聊天室喽\n"))
	scanner := bufio.NewScanner(conn)
Start1:
	for scanner.Scan() {
		name := chat.GetName(conn)
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
			chat.PrivateMessage(conn, fmt.Sprintf("[私聊] 【%s】对你说：%s \n", name, split[1]), targetName)

		} else if strings.HasPrefix(msg, "/") {
			chose := strings.TrimPrefix(msg, "/")
			switch chose {
			case "exit":
				log.Printf("用户%s输入/exit自动退出", name)
				goto FORFINISH
			case "who":
				onlineUsers := make([]string, 0)
				for _, name := range chat.clients {
					if name != "" {
						onlineUsers = append(onlineUsers, name)
					}
				}
				conn.Write([]byte("当前聊天室在线的用户列表如下" +
					fmt.Sprintf(strings.Join(onlineUsers, " & ")) + "\n"))
			case "rename":
				conn.Write([]byte("请输入您想更改的名字；" + "\n"))
				if scanner.Scan() {
					newname := scanner.Text()
					chat.Rename(conn, newname)
				}
			}
		} else {
			chat.BroadCast(conn, fmt.Sprintf("[%s]: %s\n", name, msg))
			log.Println("广播信息已经投放...")
		}

	}
FORFINISH:
	fmt.Printf("出现了退出错误")
	//sanner.Scan()如果断开，意味着用户已经断开
	chat.mu.Lock()
	delete(chat.clients, conn)
	chat.mu.Unlock()
	chat.BroadCast(conn, fmt.Sprintf("[%s] 离开了聊天室", chat.GetName(conn)))

}

// 全体广播
func (chat *ChatRoom) BroadCast(sender net.Conn, msg string) {
	chat.mu.Lock()
	conns := make([]net.Conn, 0)
	for conn := range chat.clients {
		if conn != sender {
			conns = append(conns, conn)
		}
	}
	chat.mu.Unlock()

	for _, conn := range conns {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Printf("给%s发消息失败，跳过", conn.RemoteAddr())
			continue
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
