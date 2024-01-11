package stores

import "strconv"

const (
	OP_NOP = iota
	OP_EQ
	OP_LIKE
)

type FilterInt64 struct {
	Op    int
	Value int64
}

func (f *FilterInt64) Set(id string, values map[string][]string) {
	if v, ok := values[id]; ok {
		if len(v) > 0 {
			if vi, err := strconv.ParseInt(v[0], 10, 64); err == nil {
				f.Eq(vi)
				return
			}
		}
	}
}

func (f *FilterInt64) Eq(value int64) {
	f.Op = OP_EQ
	f.Value = value
}

type FilterString struct {
	Value string
	Op    int
}

func (f *FilterString) Set(id string, values map[string][]string) {
	if v, ok := values[id]; ok {
		if len(v) > 0 {
			f.Eq(v[0])
			return
		}
	}
	if v, ok := values[id+".like"]; ok {
		if len(v) > 0 {
			f.Like(v[0])
			return
		}
	}
}

func (f *FilterString) Eq(value string) {
	f.Op = OP_EQ
	f.Value = value
}

func (f *FilterString) Like(value string) {
	f.Op = OP_LIKE
	f.Value = value
}
