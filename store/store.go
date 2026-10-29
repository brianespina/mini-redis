package store

import (
	"errors"
	"slices"
	"sync"
)

var (
	mu sync.RWMutex
)

var (
	ValueFetchTypeMissmatch = errors.New("Wrong type of value being fetched")
	KeyFetchDoesNotExist    = errors.New("The key being fetched does not exist")
	IndexOutOfRange         = errors.New("The index is out of range")
	KindErrorNotImplemented = errors.New("The value kind is not implemented")
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
			return "", ValueFetchTypeMissmatch
		}
		return v.str, nil
	} else {
		return "", KeyFetchDoesNotExist
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

func GetRange(key string, start int, end int) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	val, exist := store[key]

	if exist {
		if val.kind != KindList {
			return []string{}, ValueFetchTypeMissmatch
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
			return []string{}, IndexOutOfRange
		}

		return val.list[start : end+1], nil
	}

	return []string{}, KeyFetchDoesNotExist
}

func PrependList(key string, value ...string) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := store[key]

	if !ok {
		v = Value{kind: KindList}
	} else if v.kind != KindList {
		return 0, ValueFetchTypeMissmatch
	}

	buff := make([]string, 0, len(value)+len(v.list))
	for _, val := range slices.Backward(value) {
		buff = append(buff, val)
	}

	v.list = append(buff, v.list...)
	store[key] = v

	return len(v.list), nil

}

func AppendList(key string, value ...string) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := store[key]

	if !ok {
		v = Value{kind: KindList}
	} else if v.kind != KindList {
		return 0, ValueFetchTypeMissmatch
	}

	v.list = append(v.list, value...)
	store[key] = v

	return len(v.list), nil
}

func Exists(key string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := store[key]
	if !ok {
		return false
	}
	return true
}

func ValLen(key string) (int, error) {
	v, ok := store[key]
	if !ok {
		return 0, KeyFetchDoesNotExist
	}

	if v.kind == KindList {
		return len(v.list), nil
	}

	if v.kind == KindString {
		return len(v.str), nil
	}

	return 0, KindErrorNotImplemented
}
