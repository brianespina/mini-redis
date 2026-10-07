package main

import (
	"bufio"
	"fmt"
	"io"
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

			store.SetString(args[1], args[2])
			c.Write([]byte("+OK\r\n"))

		case "DEL":
			if len(args) < 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}
			count := store.Delete(args)
			fmt.Fprintf(c, ":%d\r\n", count)

		case "GET":
			if len(args) != 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}
			str, err := store.GetString(args[1])
			if err != nil {
				c.Write([]byte("%-1\r\n"))
				break
			}

			fmt.Fprintf(c, "+%s\r\n", str)

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

func WriteSerializedArray(c io.Writer, arr []string) {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(arr))

	for _, s := range arr {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(s), s)
	}

	c.Write([]byte(b.String()))
}
