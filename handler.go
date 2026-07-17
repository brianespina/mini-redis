package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConn(c net.Conn) {
	defer c.Close()
	fmt.Println("client connected", c.RemoteAddr())

	r := bufio.NewReader(c)

	for {

		fmt.Println("run")

		args, err := readCommands(r)
		if err != nil {
			fmt.Println("Client disconnected", err)
			return
		}

		if len(args) == 0 {
			continue
		}

		switch strings.ToUpper(args[0]) {
		case "COMMAND":
			c.Write([]byte("*0\r\n"))
		case "PING":
			c.Write([]byte("+PONG\r\n"))
		default:
			c.Write([]byte("-ERR unknown command\r\n"))
		}

	}
}
