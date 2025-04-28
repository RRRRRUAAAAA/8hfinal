package utils

import (
	"net"
	"strings"
)

func ReadMessages(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	recvMsg := string(buf[:n])
	recvMsg = strings.TrimSpace(recvMsg)
	return recvMsg, nil
}
