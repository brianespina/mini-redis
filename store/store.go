package store

import (
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

func Get(key string) (Value, bool) {
	v, ok := store[key]
	return v, ok
}

func SetString(key string, val string) {
	store[key] = Value{kind: KindString,
		str: val,
	}
}
