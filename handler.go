package main

import (
	"bufio"
	"fmt"
	"mini-redis/store"
	"net"
	"strconv"
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

			count := 0
			mu.Lock()
			for _, key := range args[1:] {
				_, exits := store[key]
				if exits {
					delete(store, key)
					count++
				}
			}
			mu.Unlock()
			fmt.Fprintf(c, ":%d\r\n", count)

		case "GET":
			if len(args) != 2 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

			mu.RLock()
			val, ok := store.Get(args[1])

			if ok {
				if val.GetValueKind() != store.KindString {
					c.Write([]byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
					break
				}
				s := val.GetStr()
				fmt.Fprintf(c, "$%d\r\n%s\r\n", len(s), s)
			} else {
				c.Write([]byte("$-1\r\n"))
			}
			mu.RUnlock()

		case "LRANGE":
			if len(args) != 4 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

			existing, exist := store[args[1]]
			if exist {
				l, ok := existing.([]string)
				if !ok {
					c.Write([]byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
					break
				}

				start, err := strconv.Atoi(args[2])
				if err != nil {
					c.Write([]byte("-value is not an integer or out of range"))
					break
				}
				end, err := strconv.Atoi(args[3])
				if err != nil {
					c.Write([]byte("-value is not an integer or out of range"))
					break
				}

				n := len(l)

				if start < 0 {
					start = n + start
				}
				if end < 0 {
					end = n + end
				}

				if start < 0 {
					start = 0
				}
				if end >= n {
					end = n - 1
				}

				if start > end || start > n {
					WriteSerializedArray(c, []string{})
					break
				}

				WriteSerializedArray(c, l[start:end+1])
				break
			}

			c.Write([]byte("*0\r\n"))
		case "LPUSH":
			if len(args) != 3 {
				c.Write([]byte("-ERR wrong number of arguments\r\n"))
				break
			}

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

		default:
			c.Write([]byte("-ERR unknown command\r\n"))
		}

	}
}

func WriteSerializedArray(c net.Conn, arr []string) {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(arr))

	for _, s := range arr {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(s), s)
	}

	c.Write([]byte(b.String()))
}
