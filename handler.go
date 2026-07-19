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
	store := make(map[string]string)

	for {
		args, err := readCommands(r)
		if err != nil {
			fmt.Println("Client disconnected", err)
			return
		}

		if len(args) == 0 {
			continue
		}
		fmt.Println(args)

		switch strings.ToUpper(args[0]) {
		case "COMMAND":
			c.Write([]byte("*0\r\n"))
		case "PING":
			c.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(args) != 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

			response := fmt.Sprintf("$%d\r\n%s\r\n", len(args[1]), args[1])
			c.Write([]byte(response))
		case "SET":
			if len(args) != 3 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}
			store[args[1]] = args[2]
			c.Write([]byte("+OK\r\n"))

		case "GET":
			if len(args) != 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}
			if val, ok := store[args[1]]; ok {
				if _, err := fmt.Fprintf(c, "$%d\r\n%s\r\n", len(val), val); err != nil {
					c.Write([]byte("-ERR writing to writer\r\n"))
					break
				}
			}
			c.Write([]byte("$-1\r\n"))
		default:
			c.Write([]byte("-ERR unknown command\r\n"))
		}

	}
}
