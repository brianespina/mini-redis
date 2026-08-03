package store

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

func Get(key string) (Value, bool) {
	v, ok := store[key]
	return v, ok
}

func (v Value) GetStr() string {
	return v.str
}

func (v Value) GetValueKind() ValueKind {
	return v.kind
}
