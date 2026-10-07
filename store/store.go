package store

import (
	"errors"
	"strconv"
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

var store = map[string]Value{}

func GetString(key string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := store[key]
	if ok {
		if v.kind != KindString {
			return "", errors.New("WRONGTYPE")
		}
		return v.str, nil
	} else {
		return "", errors.New("NONEXISTING")
	}
}

func SetString(key string, val string) {
	mu.Lock()
	defer mu.Unlock()
	store[key] = Value{kind: KindString,
		str: val,
	}
}

func Delete(args []string) int {

	count := 0
	mu.Lock()
	defer mu.Unlock()
	for _, key := range args[1:] {
		_, exist := store[key]
		if exist {
			delete(store, key)
			count++
		}
	}
	return count
}

func GetList(args []string) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	val, exist := store[args[1]]

	if exist {
		if val.kind != KindList {
			return []string{}, errors.New("WRONGTYPE")
		}

		start, err := strconv.Atoi(args[2])
		if err != nil {
			return []string{}, errors.New("NOT INT OR OUT OF RANGE")
		}

		end, err := strconv.Atoi(args[3])
		if err != nil {
			return []string{}, errors.New("NOT INT OR OUT OF RANGE")
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
			return []string{}, errors.New("NOT INT OR OUT OF RANGE")
		}

		return val.list[start : end+1], nil
	}

	return []string{}, errors.New("NONEXISTING")
}
