package channelchatroom

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Ip         string
	port       int
	Clients    map[net.Conn]*User
	ConnNumber int
	Message    chan string
	mutex      sync.Mutex
	Commands   []Command
}

// 第一个调用，先设置一下端口
func (s *Server) InitCommands() {
	s.Commands = []Command{
		{"/exit", "退出聊天室"},
		{"/who", "查看在线用户"},
		{"/rename", "修改用户名"},
		// {"/kick", "踢出其他用户"},
		// {"/mute", "禁言其他用户"},
	}
}
func SetSever() *Server {
	var ip string
	var port int
	fmt.Println("请输入服务器ip:")
	fmt.Scanln(&ip)
	fmt.Println("请输入服务器端口:")
	fmt.Scanln(&port)
	server := NewServer(ip, port)
	return server
}
func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:         ip,
		port:       port,
		Clients:    make(map[net.Conn]*User),
		Message:    make(chan string, 1024),
		ConnNumber: 0, // 防止堵塞
		Commands:   make([]Command, 0),
	}
	s.InitCommands()
	return s
}

func (s *Server) StartServer() {
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.port))
	defer Listener.Close()
	if err != nil {
		log.Println("监听器出现错误：", err)
	}
	log.Println("服务器启动")
	for {
		conn, err := Listener.Accept()
		if err != nil {
			log.Println("服务器连接出现问题：")
		}
		log.Println("连接成功")
		//新建user的过程（map中需要）
		user := s.SetUser(conn)
		s.AskServerDatil(conn, user)

		go s.HandleMsg(conn)
	}

}

// 给信息赋值
func (s *Server) AskServerDatil(conn net.Conn, user *User) {
	s.ConnNumber++
	s.mutex.Lock()
	s.Clients[conn] = user
	s.mutex.Unlock()
}

// 处理信息
func (s *Server) HandleMsg(conn net.Conn) {
	for {
		s.ReceiveMsg(conn)
	}
}

// 发送信息
func (s *Server) ReceiveMsg(conn net.Conn) {
	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	context := buf[:n]
	if string(context) != "" {
		s.ClassficationMsg(string(context), conn)
	} else {
		conn.Write([]byte("不要发送空消息\n"))
	}

}

func (s *Server) ClassficationMsg(context string, conn net.Conn) {
	if strings.HasPrefix(context, "@") {
		s.PrivateHandle(context, conn)
	} else if strings.HasPrefix(context, "/") {
		fmt.Println("命令消息")
	} else {
		fmt.Println("群聊消息")
	}
}

func (s *Server) PrivateHandle(context string, conn net.Conn) {
	successful := false
	context = strings.TrimPrefix(context, "@")
	text := strings.SplitN(context, " ", 2)
	s.mutex.Unlock()
	defer s.mutex.Unlock()
	for connList := range s.Clients {
		if s.Clients[connList].Name == text[0] {
			_, err := connList.Write([]byte(text[1]))
			if err != nil {
				log.Println("写消息出现错误")
			}
			successful = true
			break
		}
	}
	if successful {
		fmt.Printf("您的信息已经发送给：%s", text[0])
	} else {
		fmt.Printf("消息发送失败，未找到该用户")
	}
}

func (s *Server) CommandHandle(context string, conn net.Conn) {
	context = strings.TrimPrefix(context, "/")
	switch context {
	//查询所有指令
	case "all":
		s.ShowCommand()
	case "who":
		var onlineUsers []string
		for i := range s.Clients {
			onlineUsers = append(onlineUsers, s.Clients[i].Name)
		}
		fmt.Printf(strings.Join(onlineUsers, "&") + "\n")
	case "exit":
		log.Printf("用户%s断开连接", s.Clients[conn].Name)
		s.mutex.Lock()
		delete(s.Clients, conn)
		s.mutex.Unlock()
		conn.Close()
	case "rename":
		fmt.Println("请输入新名字:")
		var newname string
		oldname := s.Clients[conn].Name
		fmt.Scanln(&newname)
		if newname == oldname {
			fmt.Println("请不要输入和原来重复的名字")
		}
		for i := range s.Clients {
			if s.Clients[i].Name == newname {
				fmt.Printf("%s已存在，改名失败", newname)
			}
			s.Clients[conn].Name = newname
			fmt.Printf("改名成功！新名字为%s", newname)
		}
	}

}

// 展示所有命令
func (s *Server) ShowCommand() {
	for name, desc := range s.Commands {
		fmt.Printf("%s:%s", name, desc)
	}
}

func (s *Server) BroadcastHandle(context string, conn net.Conn) {
	fmt.Println("处理广播消息")
}

func (s *Server) SetUser(conn net.Conn) *User {
	user := NewUser()
	user.AskUserdetail(conn)
	return user
}
