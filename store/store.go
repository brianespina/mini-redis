package store

import (
	"errors"
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

func GetList(key string, start int, end int) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	val, exist := store[key]

	if exist {
		if val.kind != KindList {
			return []string{}, errors.New("WRONGTYPE")
		}

		n := len(val.list)

		if start < 0 {
			start += n
		}
		if end < 0 {
			end += n
		}

		if start < 0 {
			start = 0
		}

		if end >= n {
			end = n - 1
		}

		if start > end {
			return []string{}, errors.New("NOT INT OR OUT OF RANGE")
		}

		return val.list[start : end+1], nil
	}

	return []string{}, errors.New("NONEXISTING")
}

func SetList(key string, value ...string) (int, error) {
	return 0, nil
}
