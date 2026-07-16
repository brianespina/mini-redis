package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("listening to port")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	fmt.Println("client connected", c.RemoteAddr())

	r := bufio.NewReader(c)

	for {
		//get array len
		n, err := readInt(r, "*")
		if err != nil {
			fmt.Println("Client disconnected", err)
			return
		}

		//args buffer
		args := []string{}

		//loop array len times
		for range n {
			command, err := readCommand(r)
			if err != nil {
				fmt.Println("Client disconnected", err)
				return
			}
			args = append(args, command)
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
