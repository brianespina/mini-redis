package store

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

var (
	mu sync.RWMutex
)

type ValueKind int

const (
	KindString ValueKind = iota
	KindList
)

type Value struct {
	kind ValueKind
	str  string
	list []string
}

func (v Value) GetStr() string {
	return v.str
}

func (v Value) GetValueKind() ValueKind {
	return v.kind
}

var store = map[string]Value{}

func GetString(c io.Writer, key string) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := store[key]
	if ok {
		if v.GetValueKind() != KindString {
			c.Write([]byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
			return
		}
		s := v.GetStr()
		fmt.Fprintf(c, "$%d\r\n%s\r\n", len(s), s)
	} else {
		c.Write([]byte("$-1\r\n"))
	}
}

func SetString(c io.Writer, key string, val string) {
	mu.Lock()
	defer mu.Unlock()
	store[key] = Value{kind: KindString,
		str: val,
	}

	c.Write([]byte("+OK\r\n"))
}

func Delete(c io.Writer, args []string) {

	count := 0
	mu.Lock()
	for _, key := range args[1:] {
		_, exist := store[key]
		if exist {
			delete(store, key)
			count++
		}
	}
	mu.Unlock()
	fmt.Fprintf(c, ":%d\r\n", count)

}

func GetList(c io.Writer, args []string) {
	val, exist := store[args[1]]

	if exist {
		if val.GetValueKind() != KindList {
			c.Write([]byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
			return
		}

		start, err := strconv.Atoi(args[2])
		if err != nil {
			c.Write([]byte("-value is not an integer or out of range"))
			return
		}
		end, err := strconv.Atoi(args[3])
		if err != nil {
			c.Write([]byte("-value is not an integer or out of range"))
			return
		}

		n := len(val.list)

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
			return
		}

		WriteSerializedArray(c, val.list[start:end+1])
		return
	}

	c.Write([]byte("*0\r\n"))
}

func WriteSerializedArray(c io.Writer, arr []string) {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(arr))

	for _, s := range arr {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(s), s)
	}

	c.Write([]byte(b.String()))
}
