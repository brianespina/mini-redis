package main

type ValueKind int

const (
	KindString ValueKind = iota
	KindList
)

type Value struct {
	Kind ValueKind
	Str  string
	List []string
}

var store = make(map[string]Value)

func newStringVal(s string) Value {
	return Value{Kind: KindString, Str: s}
}
