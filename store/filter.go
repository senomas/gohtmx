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

func (f *FilterInt64) EQ(value int64) {
	f.Op = OP_EQ
	f.Value = value
}

type FilterString struct {
	Op    int
	Value string
}

func (f *FilterString) EQ(value string) {
	f.Op = OP_EQ
	f.Value = value
}

func (f *FilterString) LIKE(value string) {
	f.Op = OP_LIKE
	f.Value = value
}
