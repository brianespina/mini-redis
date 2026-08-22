package main

import (
	"bufio"
	"fmt"
	"mini-redis/store"
	"net"
	"strings"
)

func handleConn(c net.Conn) {
	defer c.Close()
	fmt.Println("client connected", c.RemoteAddr())

	r := bufio.NewReader(c)

	for {
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
			WriteSerializedArray(c, []string{})
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
			store.SetString(c, args[1], args[2])

		case "DEL":

			if len(args) < 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

			store.Delete(c, args)

		case "GET":
			if len(args) != 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}
			store.GetString(c, args[1])

		case "LRANGE":
			if len(args) != 4 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

		case "LPUSH":
			if len(args) != 3 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

			/*
				mu.Lock()
				existing, exist := store[args[1]]

				var list []string

				if exist {
					l, ok := existing.([]string)
					if !ok {
						mu.Unlock()
						c.Write([]byte("-ERR trying to push to a non list\r\n"))
						break
					}
					list = l
				}

				list = append(list, args[2])
				store[args[1]] = list
				fmt.Fprintf(c, ":%d\r\n", len(list))
				mu.Unlock()
			*/
		default:
			c.Write([]byte("-ERR unknown command\r\n"))
		}

	}
}
