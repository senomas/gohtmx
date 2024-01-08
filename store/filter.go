package store

const (
	OP_NOP = iota
	OP_EQ
	OP_LIKE
)

type FilterInt64 struct {
	Op    int
	Value int64
}

func (f *FilterInt64) Eq(value int64) {
	f.Op = OP_EQ
	f.Value = value
}

type FilterString struct {
	Value string
	Op    int
}

func (f *FilterString) Eq(value string) {
	f.Op = OP_EQ
	f.Value = value
}

func (f *FilterString) Like(value string) {
	f.Op = OP_LIKE
	f.Value = value
}
